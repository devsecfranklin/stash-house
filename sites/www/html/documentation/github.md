# x_badge

- fork the dead10c5 version of the repo
  - a copy will show up in your personal github
- make a clone of your copy of the repo
- set the origin repo as your upstream from the cli

```sh
git remote add upstream git@github.com:DEAD10C5/x_badge.git
```

- Now when you do a pull request it will go to the original repo.
- Make a new branch

```sh
git checkout -b dev1
```

- Make changes
- Commit changes

```sh
git commit -m 'I did something'
```

- push changes

```sh
git push origin dev1
```

- Now go to the github website and you should see a big yellow bar with a button
  to open a pull reqest.
  - have to be at top level of repo
- add a description of changes and submit the PR
- now we can review and comment
