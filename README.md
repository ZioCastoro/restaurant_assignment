# restaurant_assignment
Restaurant Assignment

To launch the project use the command `make all`.

To list commands, use the command `make help`.

Project folders:

- app: application code.
    - entities: entity objects.
    - handlers: code to handle http requests and responses.
        - payloads: requests and responses objects and mapping.
    - repositories: code to handle interation with data substrate.
    - services: code containing application logic.
    - validators: validation library.
- cmd: application main.
- deps: dependency injection.
- json: json encoding utilities.
- postgre: postgreSQL library.
- scripts: database migration scripts.
- utils: miscellaneous utilities.