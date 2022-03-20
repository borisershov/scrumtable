# ScrumTable Bot

[ScrumTable Bot](https://t.me/scrumtable_bot) represents a simple tool to orgonize your tasks. With this bot in Telegram you are able to create a short issues you need to complete for a day or a week (sprint).

It is completely free application you can either use a [ready-made bot](https://t.me/scrumtable_bot) or built and deploy it in your own infrastructure.

## Overview

### Calendar

Bot provides you a calendar for select a date to watch your already created issues or make a new one. With this calendar you may choose two types of dates:
- Sprint (first column in dates table means a number of week in the year; starts with `#` character)
- Day (rest columns in dates table)

To see the calendar execute a `/calendar` command in a bot menu.

### Sprint issues

In [ScrumTable Bot](https://t.me/scrumtable_bot) the sprint is equals to one week long. Sprint issues are issues you need to complete within a given week. The one of a sprint issues must be a `goal`. You may have only one goal issue within a sprint. First created sprint issue automatically will set as `goal`.

To create a new sprint issue you only need to send a text within an appropriate bot menu.

You may click for each sprint issue to do following actions with them:
- Mark complete issues as `done`
- Set issue as `goal`
- `Edit` issue text
- `Delete` issue

### Daily issues

Daily issues are issues you need to complete within a given day.

To create a new daily issue you only need to send a text within an appropriate bot menu.

You may click for each daily issue to do following actions with them:
- Mark complete issue as `done`
- Change an issue `due date`
- `Edit` issue text
- `Delete` issue

## Developers info

[ScrumTable Bot](https://t.me/scrumtable_bot) bases on a following Nixys libraries:
- [nxs-go-appctx](https://github.com/nixys/nxs-go-appctx)
- [nxs-go-conf](https://github.com/nixys/nxs-go-conf)
- [nxs-go-telegram](https://github.com/nixys/nxs-go-telegram)
