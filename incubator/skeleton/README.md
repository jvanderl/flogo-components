# Skeletons

In this folder you can find skeletons for creating flogo activities and triggers from scratch.

The idea is to copy the whole folder to your own repository of flogo componentes, rename it to your package name and replace all stuff between brackets in all the files in that folder.
Below per type the available items for each type.

## Activity

The skeleton for an activity can be found [here](activity).

| Item | Description | Example |
|:-----|:------------|:--------|
| [package] | Package name, this will be the short name used in all references| replace |
| [description] | Package description | Search and Replace |
| [author_name] | Author of this package, you ;) | Jan van der Lugt |
| [author_email] | Email address of the author of this package) | jvanderl@tibco.com |
| [functionality] | functionality of the activity | search and replace strings |
| [git_user] | Git user name | jvanderl
| [git_repo] | Git reposotory | flogo-components |
| [input(x)_name] | Name for Input x (use as many as you need) | searchstring |
| [input(x)_type] | Type for Input x | string, int, bool etc. |
| [input(x)_default] | Default value for Input x | replace me |
| [input(x)_value] | Value for Input x used only for example code | : |
| [input(x)_desc] | Description for Input x | The string to be replaced
| [output(x)_name] | Name for Output x (use as many as you need) | result |
| [ouput(x)_type] | Type for Output x | string, int, bool etc. |
| [import(x)] | Imported library, will end up in your code (use as many as you need) | strings

## Trigger
**- Under Construction -**