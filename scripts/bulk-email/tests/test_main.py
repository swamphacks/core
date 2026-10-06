import csv

import pytest

import main


class FakeSES:
    def __init__(self):
        self.send_calls = []

    def send_bulk_templated_email(self, **kwargs):
        self.send_calls.append(kwargs)
        return {"Status": [{"Status": "Success"}]}


def write_contacts(path, count):
    with path.open("w", newline="") as file:
        writer = csv.DictWriter(file, fieldnames=["email"])
        writer.writeheader()
        for index in range(count):
            writer.writerow({"email": f"person{index}@example.com"})


def test_parse_csv_reads_contacts(tmp_path):
    csv_path = tmp_path / "contacts.csv"
    write_contacts(csv_path, 2)

    contacts = main.parse_csv(csv_path)

    assert [contact.get_email() for contact in contacts] == [
        "person0@example.com",
        "person1@example.com",
    ]


def test_parse_csv_rejects_missing_email(tmp_path):
    csv_path = tmp_path / "contacts.csv"
    csv_path.write_text("name\nNam\n")

    with pytest.raises(ValueError, match="Missing email field"):
        main.parse_csv(csv_path)


def test_chunk_list_splits_at_requested_size():
    assert list(main.chunk_list(list(range(105)), 50)) == [
        list(range(50)),
        list(range(50, 100)),
        list(range(100, 105)),
    ]


def test_send_emails_batches_without_real_aws(tmp_path, monkeypatch):
    csv_path = tmp_path / "contacts.csv"
    write_contacts(csv_path, 51)
    fake_ses = FakeSES()
    monkeypatch.setattr(main, "get_ses_client", lambda: fake_ses)

    main.send_emails("test-template", csv_path)

    assert len(fake_ses.send_calls) == 2
    assert len(fake_ses.send_calls[0]["Destinations"]) == 50
    assert len(fake_ses.send_calls[1]["Destinations"]) == 1
    assert fake_ses.send_calls[0]["Template"] == "test-template"
    assert fake_ses.send_calls[0]["Source"] == (
        "SwampHacks <contact@swamphacks.com>"
    )


def test_missing_template_does_not_create_aws_client(
    tmp_path, monkeypatch, capsys
):
    monkeypatch.chdir(tmp_path)

    def fail_if_called():
        raise AssertionError("AWS client should not be created")

    monkeypatch.setattr(main, "get_ses_client", fail_if_called)

    main.upload_template("test-template", "missing.html")

    assert "File not found" in capsys.readouterr().out
