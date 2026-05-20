#!/usr/bin/env python3
"""Interactive client for creating documents through the PTMK service."""

from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.parse import urljoin
from urllib.request import Request, urlopen


DEFAULT_BASE_URL = "http://localhost:8080"
DOCUMENTS_PATH = "/documents"


def prompt_text(label: str, *, required: bool = True, default: str | None = None) -> str | None:
    prompt = label
    if default is not None:
        prompt += f" [{default}]"
    prompt += ": "

    while True:
        value = input(prompt).strip()
        if value:
            return value
        if default is not None:
            return default
        if not required:
            return None
        print("Поле обязательно, введите значение.")


def prompt_int(label: str, *, default: int | None = None) -> int:
    while True:
        raw = prompt_text(label, default=str(default) if default is not None else None)
        try:
            value = int(raw or "")
        except ValueError:
            print("Введите целое число.")
            continue

        if value < 1:
            print("Введите число больше 0.")
            continue
        return value


def prompt_date(label: str, *, required: bool = True) -> str | None:
    while True:
        value = prompt_text(label, required=required)
        if value is None:
            return None
        try:
            datetime.strptime(value, "%Y-%m-%d")
        except ValueError:
            print("Дата должна быть в формате YYYY-MM-DD, например 2026-03-01.")
            continue
        return value


def prompt_choice(label: str, choices: dict[str, str]) -> str:
    options = ", ".join(f"{key} - {value}" for key, value in choices.items())
    while True:
        value = prompt_text(f"{label} ({options})")
        if value in choices:
            return choices[value]
        print("Выберите один из доступных вариантов.")


def prompt_yes_no(label: str, *, default: bool = False) -> bool:
    default_hint = "Y/n" if default else "y/N"
    while True:
        value = input(f"{label} [{default_hint}]: ").strip().lower()
        if not value:
            return default
        if value in {"y", "yes", "д", "да"}:
            return True
        if value in {"n", "no", "н", "нет"}:
            return False
        print("Введите y/yes/да или n/no/нет.")


def prompt_legal_entity() -> dict[str, Any]:
    legal_entity: dict[str, Any] = {
        "name": prompt_text("Название юридического лица"),
    }

    entity_type = prompt_text("Тип юрлица", required=False)
    if entity_type is not None:
        legal_entity["entity_type"] = entity_type

    return {
        "type": "LEGAL_ENTITY",
        "legal_entity": legal_entity,
    }


def prompt_natural_person() -> dict[str, Any]:
    while True:
        natural_person = {
            "first_name": prompt_text("Имя", required=False),
            "middle_name": prompt_text("Отчество", required=False),
            "last_name": prompt_text("Фамилия", required=False),
            "initials": prompt_text("Инициалы", required=False),
        }

        if any(value is not None for value in natural_person.values()):
            return {
                "type": "NATURAL_PERSON",
                "natural_person": natural_person,
            }

        print("Для физического лица нужно заполнить хотя бы одно поле имени.")


def prompt_persons() -> list[dict[str, Any]]:
    count = prompt_int("Количество участников документа", default=1)
    persons = []

    for index in range(1, count + 1):
        print(f"\nУчастник {index}")
        person_type = prompt_choice(
            "Тип участника",
            {
                "1": "LEGAL_ENTITY",
                "2": "NATURAL_PERSON",
            },
        )
        if person_type == "LEGAL_ENTITY":
            persons.append(prompt_legal_entity())
        else:
            persons.append(prompt_natural_person())

    return persons


def build_payload() -> dict[str, Any]:
    print("Введите данные документа. Даты вводятся в формате YYYY-MM-DD.\n")

    return {
        "user_id": prompt_int("ID пользователя"),
        "document_type": prompt_text("Тип документа"),
        "document_number": prompt_text("Номер документа"),
        "date_created": prompt_date("Дата подписания/создания"),
        "date_end": prompt_date("Дата истечения", required=False),
        "persons": prompt_persons(),
    }


def post_json(base_url: str, payload: dict[str, Any]) -> tuple[int, str]:
    endpoint = urljoin(base_url.rstrip("/") + "/", DOCUMENTS_PATH.lstrip("/"))
    body = json.dumps(payload, ensure_ascii=False).encode("utf-8")
    request = Request(
        endpoint,
        data=body,
        headers={
            "Content-Type": "application/json",
            "Accept": "application/json",
        },
        method="POST",
    )

    with urlopen(request, timeout=15) as response:
        response_body = response.read().decode("utf-8")
        return response.status, response_body


def print_response(status: int, body: str) -> None:
    print(f"\nОтвет сервиса: HTTP {status}")
    if not body:
        return

    try:
        parsed = json.loads(body)
    except json.JSONDecodeError:
        print(body)
        return

    print(json.dumps(parsed, ensure_ascii=False, indent=2))


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Собрать JSON документа из пользовательского ввода и отправить его в PTMK service.",
    )
    parser.add_argument(
        "--url",
        default=None,
        help=f"Базовый адрес сервиса. По умолчанию будет предложен {DEFAULT_BASE_URL}.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Только показать JSON, без отправки в сервис.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    base_url = args.url or prompt_text("Адрес сервиса", default=DEFAULT_BASE_URL)
    payload = build_payload()

    print("\nJSON для отправки:")
    print(json.dumps(payload, ensure_ascii=False, indent=2))

    if args.dry_run or not prompt_yes_no("\nОтправить JSON в сервис?", default=True):
        print("Отправка отменена.")
        return 0

    try:
        status, body = post_json(str(base_url), payload)
    except HTTPError as error:
        error_body = error.read().decode("utf-8")
        print_response(error.code, error_body)
        return 1
    except URLError as error:
        print(f"Не удалось подключиться к сервису: {error.reason}", file=sys.stderr)
        return 1
    except TimeoutError:
        print("Сервис не ответил за отведенное время.", file=sys.stderr)
        return 1

    print_response(status, body)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
