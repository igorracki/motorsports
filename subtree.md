## Repo is configured with a remote subtree for the fronted.

```
git remote add frontend_remote git@github.com:igorracki/v0-motorsport-event-landing-page.git
git subtree add --prefix frontend frontend_remote main --squash
```

### To remove the subtree

```
git rm -r frontend/
git remote remove frontend_remote
git commit -m 'removed subtree'
```

### To update the subtree remote

```
git remote set-url frontend_remote https://github.com/username/new_frontend.git
git subtree pull --prefix frontend frontend_remote main --squash
```

### Pushing to the subtree from here

```
git subtree push --prefix frontend frontend_remote main
```
