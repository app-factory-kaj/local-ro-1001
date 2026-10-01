Feature: Todo management

  @story-1
  Rule: The todo list shows every todo in creation order, oldest first

    Scenario: Viewing todos added in sequence
      Given a User has added "Buy groceries" and then "Read a book"
      When the User views the todo list
      Then "Buy groceries" appears before "Read a book" in the list

    Scenario: Viewing an empty list
      Given no todos have been added
      When the User views the todo list
      Then the list is empty

  @story-2
  Rule: A new todo is added with a title and starts as not done

    Scenario: Adding a todo with a title
      Given no todos have been added
      When the User adds a todo titled "Walk the dog"
      Then the list has exactly one todo
      And that todo is titled "Walk the dog" and is not done

    @negative
    Scenario: A blank title is refused
      Given no todos have been added
      When the User tries to add a todo with a blank title
      Then the list is still empty

  @story-3
  Rule: Toggling a todo flips its done status

    Scenario: Marking a todo done
      Given a User has added "Walk the dog"
      When the User toggles "Walk the dog"
      Then "Walk the dog" is done

    Scenario: Marking a done todo not done again
      Given a User has added "Walk the dog" and toggled it done
      When the User toggles "Walk the dog" again
      Then "Walk the dog" is not done

  @story-4
  Rule: Deleting a todo removes it from the list

    Scenario: Removing a todo
      Given a User has added "Walk the dog" and "Read a book"
      When the User deletes "Walk the dog"
      Then the list has exactly one todo
      And it is titled "Read a book"
