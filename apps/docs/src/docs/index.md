# Docs Overview

Documentation system to make sure everyone knows how everything works. Vital for the longevity of any long term software engineering project.

## Rationale

We decided to make our documentation using a MkDocs project living in the same GitHub repository as everything else. Why this over something like a GitHub wiki? The answer is simple:

- Developers can add documentation changes in the same Git commits as their work. This ensures that when a new feature is added, it will be apparent whether or not related documentation was also made when a pull request is reviewed.

- GitHub Actions can hook up pushed changes to update a website, meaning no one person has to manage changes to the documentation in a seperate repository.

- Writing in markdown files can be moved to a different documentation site generator if needed, and can be easily edited by non-engineers.


