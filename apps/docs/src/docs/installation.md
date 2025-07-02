# Gettings Started

## Normal setup

1. Install [mkdocs](https://www.mkdocs.org/user-guide/installation/)
1. Install the python packages listed in core/apps/docs/requirements.txt
1. Run:

```bash
mkdocs serve
```

!!! warning
    If the site does not generate because there is still missing dependancies, you can find them using [mkdocs-get-deps](https://github.com/mkdocs/get-deps)

## For linux distributions which do not package python packages through pip

This includes Arch Linux.

1. Run
```bash
python -m venv env
source env/bin/activate
```
1. Install packages
```bash
pip install -r requirements.txt
```
1. Run with
```bash
mkdocs serve
```

---

Now just make normal changes to the documentation and the site will update automatically. Make commits and push your changes when you're done!
