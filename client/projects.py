#!/usr/bin/env python3
import time
import clapi

from gql import gql


def create_project(client):
    mutation = gql(
    """
    mutation CreateProject($np: NewProject!) {
        createProject(input: $np) {
            id
            name
            code
        }
    }
    """    
    )

    stamp = time.time()

    params = {
        "np": {
            "name": f"Project{stamp}",
            "code": f"p{stamp}"
        }
    }

    result = client.execute(mutation, variable_values=params)
    print(result)

def projects(client):
    query = gql(
    """
    query projects {
        projects {
            id
            name
            code
        }
    }
    """
    )

    result = client.execute(query)
    projects = result['projects']
    for project in projects:
        print(f"{project}")


def main():
    config = clapi.load_config()
    client = clapi.new_client(config)
    
    create_project(client)
    projects(client)


if __name__ == "__main__":
    main()


