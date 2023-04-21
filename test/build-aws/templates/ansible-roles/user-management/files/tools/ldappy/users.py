import io, re
import random
from collections import defaultdict
from jinja2 import PackageLoader, Environment


_attribute_re = re.compile(f'([^:]+):\s+(.*)')


def check_uid_number(conn, search_base, uid_number):
    entries = conn.extend.standard.paged_search(
        search_base,
        f'(&(objectClass=posixAccount)(uidNumber={uid_number}))',
        attributes=['cn', 'uidNumber'],
        paged_size=1
    )
    for entry in entries:
        return False
    return True
    

def create_user(conn, domain_dn, domain_dns,
                    user_name, user_password, 
                    given_name, family_name, 
                    uid_number, gid_number):

    # check if the user already exists
    entries = conn.extend.standard.paged_search(
        domain_dn,
        f'(&(objectClass=posixAccount)(cn={user_name}))',
        attributes=['cn', 'uidNumber'],
        paged_size=1
    )
    for entry in entries:
        raise ValueError(f"user {user_name} already exists")

    # check if the uid already exists
    entries = conn.extend.standard.paged_search(
        domain_dn,
        f'(&(objectClass=posixAccount)(uidNumber={uid_number}))',
        attributes=['cn', 'uidNumber'],
        paged_size=1
    )
    for entry in entries:
        raise ValueError(f"user {uid_number} already exists")

    # expand the user ldif template
    env = Environment(loader=PackageLoader("ldappy", "templates"))
    template = env.get_template("user.ldif")
    doc = template.render(
        domain_dn=domain_dn,
        domain_dns=domain_dns,
        user_name=user_name,
        user_password=user_password,
        given_name=given_name,
        family_name=family_name,
        uid_number=uid_number,
        gid_number=gid_number
    )
    
    print(doc)

    # parse the ldif doc
    dn = None
    oclass = []
    attrs = defaultdict(list)

    for line in io.StringIO(doc):
        mo = _attribute_re.match(line)
        if mo is None:
            raise ValueError(f"couldn't match {line}")

        key, value = mo.group(1).strip().lower(), mo.group(2).strip()

        if key == "dn":
            dn = value
            continue
        if key == "objectclass":
            oclass.append(value)
            continue
        attrs[key].append(value)

    rv = conn.add(dn, oclass, attrs)
    if rv == False:
        raise ValueError(f"failed to create user: {conn.result}")
