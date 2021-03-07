# MakeMeA

MakeMeA is a tool for GameMasters (GMs) of TableTopRPGs (TTRPGs) that allows for generating data from random tables. Markdoen and github style tables are used to create and organize the tables. Since thie README is a markdown file and contains tables, all of these examples can be used with the tool.

## Tables

MakeMeA supports several types of tables

### Lookup Table

A lookup table is the most simple of tables. The table has a name and every item on the table has the same probability of being selected. Try it with `makemea "makemea/tables/lookup table/race"`

| Race   |
| ------ |
| Elf    |
| Human  |
| Dwarf  |
| Orc    |
| Gnome  |
| Troll  |
| Goblin |

### Dice Table

A Dice table has different probabilities for each item on the table. The number and type of dice are given as part of the table and the results of the roll are used to select an item from the table. Try it with `makemea "makemea/tables/dice table/treasure`

| Treasure | 1d6 |
| -------- | --- |
| Copper   | 1-3 |
| Silver   | 4-5 |
| Gold     | 6   |

## Organizing

Every table has a name. This name is used to tell MakeMeA which table to roll on. When you have a lot of tables, organization is key. Makemea will search the current folder and subfolders for markdown files and attempt to convert any markdown tables found into tables to roll on. MakeMeA uses headers to nest tables to allow for tables to be organized. When a table is placed under a header, that header is prefixed to the table name with a "/". When a table is nested under headers and subheaders, all the headers are combined with the table name. This lets you group related tables under a header to make them easier to find. For intsance the following table can be located with the name: `makemea/organizing/weapons`

| Weapons |
| ------- |
| Sword   |
| Dagger  |
| Longbow |
| Mace    |
| Spear   |

You can see all of the tables that MakeMeA has detected by using the `--list` command. Sometimes, you'll have a bunch of subtables that are used by a parent table. If the subtables aren't meant to be used on their own, you can hide them from the listing view by italicizing the name of the table. Notice that while the below table doesn't show up under the `--list` command, it is still accessable via `makemea makemea/organizing/hidden`

| _hidden_ |
| -------- |
| Secret   |
| Mystery  |
| Illusion |

## Templates

There are a few template functions that can be used to allow for more complex table behavior. Under the hood, golang templates are used. The syntax will be familiar to go programmers but is easy enough for anyone to follow.

### lookup

The `lookup` function can be used to get a result from another table an use it as part of a different result. This lets you reuse tables in more than one place and have complex lookup results. For instance, if we wanted to have fancier versions of the weapons above, we could do the following. Try it wih: `makemea makemea/templates/lookup/fancy`

| Fancy                                                    |
| -------------------------------------------------------- |
| Shiny {{lookup "makemea/organizing/weapons" }}           |
| Glowing {{lookup "makemea/tables/dice table/treasure" }} |
| Large {{lookup "makemea/tables/lookup table/race" }}     |

When you have a large hierarchy of deeply-nested tables, it can be cumbersome to provide the full path to every table. You can use relative paths to shorten the call to lookup. The below table has both the full path to the fancy table as well as relative paths to the same table. Try it with: `makemea makemea/templates/lookup/fancier`

| Fancier                                           |
| ------------------------------------------------- |
| Rusty {{lookup "makemea/templates/lookup/fancy"}} |
| Glittering {{ lookup "./fancy"}}                  |
| Sparkling {{ lookup "./fancy"}}                   |

### roll

The `roll` function is used to roll a set of dice as part of the final result. This is great for treasure if you want to generate a random amount of some currency. Try it with `makemea makemea/templates/roll/horde`

| Horde                        |
| ---------------------------- |
| {{roll "10d100+500"}} Copper |
| {{roll "5d20+50"}} Silver    |
| {{roll "5d8+10"}} Gold       |
| {{roll "3d6"}} Platinum      |
