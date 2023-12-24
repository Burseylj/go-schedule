# go-schedule
Proof of concept scheduling app using Go, HTMX/Hyperscript and Tailwind. The goal was to make this skeleton of a scheduling app work with a traditional server rendered architecture, using HTMX to partially reload parts of the app.

https://github.com/Burseylj/go-schedule/assets/24847234/9e5990e1-df9c-4ec6-b1e3-53b6ad563990

## Running Locally

- Start the Application:
  ```bash
  air
  ```

  then go to localhost:8080/schedule
  Events can be added or delted by hitting enter
  

## Next steps
This is fairly rudimentary. I've learned that Go's templating library is not ideal for building this kind of app (interpolation is awkwar

https://github.com/Burseylj/go-schedule/assets/24847234/a57ff663-ea56-4cb1-ac8c-b78d589dc506

d). If I go much further I would probably rewrite with a more robust templating system. Grouping rows by employee team, sorting, adding collapsable and expandable rows, adding vaildation on submit are the tough features
 - Sorting would probably be done with either a full page reload, or at least a partial reload targeting the whole calendar. This is where Go's template fragments become frustrating because passing inputs between nested templates is not as comfortable, declarative or typesafe as I would like
 - HTMX has a system for vaildation but I'd want to use custom modals, or maybe a toast indicator.




