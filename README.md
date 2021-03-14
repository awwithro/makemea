# MakeMeA

MakeMeA is a tool for GameMasters (GMs) of TableTopRPGs (TTRPGs) that allows for generating data from random tables. Markdown and github style tables are used to create and organize the tables. Since this README is a markdown file and contains tables, all of these examples can be used with the tool.

## Installing

`go get -u github.com/awwithro/makemea`

## Tables

MakeMeA supports several types of tables

### Lookup Table

A lookup table is the most simple of tables. The table has a name and every item on the table has the same probability of being selected. Try it with `makemea "makemea/tables/lookuptable/race"`

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

A Dice table has different probabilities for each item on the table. The number and type of dice are given as part of the table and the results of the roll are used to select an item from the table. Try it with `makemea "makemea/tables/dicetable/treasure`

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

You can see all of the tables that MakeMeA has detected by using the `list` command. Try it with: `makemea list`

Sometimes, you'll have a bunch of subtables that are used by a parent table. If the subtables aren't meant to be used on their own, you can hide them from the listing view by italicizing the name of the table. Notice that while the below table doesn't show up under the `list` command, it is still accessable via `makemea makemea/organizing/hidden`

| _hidden_ |
| -------- |
| Secret   |
| Mystery  |
| Illusion |

## Templates

There are a few template functions that can be used to allow for more complex table behavior. Under the hood, golang templates are used. The syntax will be familiar to go programmers but is easy enough for anyone to follow. It also allows for the use of conditionals, loops, and other templating functions.

In addition the below functions, sprig template functions can be used as well. See [here](http://masterminds.github.io/sprig/) for their docs

### lookup

The `lookup` function can be used to get a result from another table an use it as part of a different result. This lets you reuse tables in more than one place and have complex lookup results. For instance, if we wanted to have fancier versions of the weapons above, we could do the following. Try it wih: `makemea makemea/templates/lookup/fancy`

| Fancy                                                   |
| ------------------------------------------------------- |
| Shiny {{lookup "makemea/organizing/weapons" }}          |
| Glowing {{lookup "makemea/tables/dicetable/treasure" }} |
| Large {{lookup "makemea/tables/lookuptable/race" }}     |

When you have a large hierarchy of deeply-nested tables, it can be cumbersome to provide the full path to every table. You can use relative paths to shorten the call to lookup. The below table has both the full path to the fancy table as well as relative paths to the same table. Try it with: `makemea makemea/templates/lookup/fancier`

| Fancier                                           |
| ------------------------------------------------- |
| Rusty {{lookup "makemea/templates/lookup/fancy"}} |
| Glittering {{ lookup "./fancy"}}                  |
| Sparkling {{ lookup "./fancy"}}                   |

`lookup` also allows an for an optional argument for performing the same lookup multiple times. The following table will result in three itmes being selected from the `fancier` table: `makemea makemea/templates/lookup/count`

| Count                    |
| ------------------------ |
| {{lookup "./fancier" 3}} |

### roll

The `roll` function is used to roll a set of dice as part of the final result. This is great for treasure if you want to generate a random amount of some currency. Try it with `makemea makemea/templates/roll/horde`

| Horde                        |
| ---------------------------- |
| {{roll "10d100+500"}} Copper |
| {{roll "5d20+50"}} Silver    |
| {{roll "5d8+10"}} Gold       |
| {{roll "3d6"}} Platinum      |

### fudge

The `fudge` function works similar to the `lookup` function but allows you to provide an alternate set of dice to roll. This is useful if you want to reuse an existing table but only want to use a subset of the times on that table. The following will roll on the treasure table put with a die range that will only allow for the silver and gold values to be rolled. Try it with: `makemea makemea/templates/fudge/goldorsilver`

| Gold or Silver                                        |
| ----------------------------------------------------- |
| {{fudge "makemea/tables/dicetable/treasure" "1d3+3"}} |

### Combining Templates

`roll` and `lookup` can be combined useing variables to lookup a value from another table a random number of times. The following table does the following:

1. Rolls 2d4
2. Stores the result in a variable named "r"
3. Calls the lookup function and uses the value of "r" to perform the lookup 2 - 8

Try it with `makemea "makemea/templates/combiningtemplates/encounter"`

| Encounter                                                           |
| ------------------------------------------------------------------- |
| {{$r := roll "2d4"}}{{lookup "makemea/tables/lookuptable/race" $r}} |

## Variables

More complex lookups can be done by using a variable to lookup a table based on the result of another lookup. In this example, an npc is generated on a single table. This is done by:

1. Doing a lookup of the Race table
2. Storing the result of the lookup into the variabls `$r`
3. Using the `$r` variable to determine which name table to use
4. Printing out an NPC using both the race and name generated

Admittedly, fitting this logic into a single table cell is a bit cumbersome. Luckily there is also an option for creating text blocks that are easily formatted and clearer to read.

Try it with `./makemea "makemea/variables/npc"`

| NPC                                                                            |
| ------------------------------------------------------------------------------ |
| {{$r:=lookup "./race" }}Race: {{$r}} Name: {{lookup (print "./" $r "/names")}} |

| _Race_ |
| ------ |
| Elven  |
| Human  |

### Elven

| _Names_ |
| ------- |
| Alluin  |
| Arwen   |
| Aegnor  |

### Human

| _Names_ |
| ------- |
| Beorn   |
| Aldor   |
| Fulgar  |

## Text

Its not quite a table but sometimes you want to generate something that perfoms lookups on other tables. Something like an NPC. It would be cumbersome to stuff everything into a table sell. Instead you can use a fenced code block. Here, the npc example from above has been redone using a code block. The result is much easier to undersatnd and maintain. Try it with `makemea makemea/text/npc`

``` npc
{{$r:=lookup "makemea/variables/race" -}}
{{$name := lookup (print  "makemea/variables/" $r "/names") -}}
{{$level := roll "2d4" -}}
Name: {{$name}}
Race: {{$r}}
Level: {{$level}}
HP: {{roll (print $level "d6")}}
```
