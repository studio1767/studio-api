package ldapgroups

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"github.com/studio1767/studio-api/internal/auth"
	"github.com/studio1767/studio-api/internal/config"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrGroupNotFound = errors.New("group not found")
)

func NewClient(cfg *config.Config, tlsCfg *tls.Config) (auth.GroupGetter, error) {

	// create the struct
	ldp := ldapClient{
		serverUri:  cfg.Ldap.ServerURI,
		tlsConfig:  tlsCfg,
		searchBase: cfg.Ldap.SearchBase,
		bindDn:     cfg.Ldap.BindDN,
		bindPw:     cfg.Ldap.BindPW,
		startTls:   cfg.Ldap.StartTLS,
		conn:       nil,
	}

	// do a test connection just to make sure it's all ok
	err := ldp.connect()
	ldp.close()
	if err != nil {
		return nil, err
	}

	return &ldp, nil
}

type ldapClient struct {
	serverUri  string
	tlsConfig  *tls.Config
	searchBase string
	bindDn     string
	bindPw     string
	startTls   bool
	conn       *ldap.Conn
}

func (ldp *ldapClient) connect() error {
	if ldp.conn != nil {
		return nil
	}

	var err error = nil
	var conn *ldap.Conn = nil
	var ldaps bool = false

	if strings.HasPrefix(ldp.serverUri, "ldaps://") {
		ldaps = true

		// remove the scheme
		server, _ := strings.CutPrefix(ldp.serverUri, "ldaps://")

		// if there's no port specified, add the default for ldaps '636'
		if strings.IndexByte(server, ':') == -1 {
			server = strings.Join([]string{server, "636"}, ":")
		}

		// dial the server
		conn, err = ldap.DialTLS("tcp", server, ldp.tlsConfig)
		if err != nil {
			return err
		}

	} else {
		// remove the scheme
		server, _ := strings.CutPrefix(ldp.serverUri, "ldap://")

		// if there's no port specified, add the default for ldap '389'
		if strings.IndexByte(server, ':') == -1 {
			server = strings.Join([]string{server, "389"}, ":")
		}

		conn, err = ldap.Dial("tcp", ldp.serverUri)
		if err != nil {
			return err
		}
	}

	// start tls if specified... unless the schema is ldaps in which case
	//   it's not needed
	if ldp.startTls && ldaps == false {
		err = conn.StartTLS(ldp.tlsConfig)
		if err != nil {
			conn.Close()
			return err
		}
	}

	// bind with the search user
	err = conn.Bind(ldp.bindDn, ldp.bindPw)
	if err != nil {
		conn.Close()
		return err
	}

	// store the connection
	ldp.conn = conn

	return nil
}

func (ldp *ldapClient) reconnect() error {
	ldp.close()
	return ldp.connect()
}

func (ldp *ldapClient) close() {
	if ldp.conn != nil {
		ldp.conn.Close()
		ldp.conn = nil
	}
}

func (ldp *ldapClient) GroupsForUser(user string) (map[string]bool, error) {
	err := ldp.connect()
	if err != nil {
		return nil, err
	}

	// if we have an email address, need to find the matching user to get the username
	if strings.IndexByte(user, '@') != -1 {
		user, err = ldp.UserNameForEmail(user)
		if err != nil {
			return nil, err
		}
	}

	// search for the groups
	searchRequest := ldap.NewSearchRequest(
		ldp.searchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", user),
		[]string{"cn"},
		nil,
	)

	sr, err := ldp.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := make(map[string]bool)
	for _, entry := range sr.Entries {
		for _, attr := range entry.Attributes {
			if attr.Name == "cn" {
				groups[attr.Values[0]] = true
				break
			}
		}
	}

	return groups, nil
}

func (ldp *ldapClient) UserNameForEmail(email string) (string, error) {
	err := ldp.connect()
	if err != nil {
		return "", err
	}

	// search for the user based on their email address
	searchRequest := ldap.NewSearchRequest(
		ldp.searchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=posixAccount)(mail=%s))", email),
		[]string{"uid"},
		nil,
	)

	sr, err := ldp.conn.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) == 1 {
		for _, attr := range sr.Entries[0].Attributes {
			if attr.Name == "uid" {
				return attr.Values[0], nil
			}
		}
	}

	return "", fmt.Errorf("%s: %w", email, ErrUserNotFound)
}

// func (ldp *ldapClient) findUser(userName, userPw string) (*User, error) {
//     l, err := ldap.Dial("tcp", ldp.server)
//     if err != nil {
//         return nil, err
//     }
//     defer l.Close()
//
//     // upgrade to tls
//     err = l.StartTLS(ldp.tlsConfig)
//     if err != nil {
//         return nil, err
//     }
//
//     // bind with the search user
//     err = l.Bind(ldp.bindDn, ldp.bindPw)
//     if err != nil {
//         return nil, err
//     }
//
//     // search for the user
//     searchRequest := ldap.NewSearchRequest(
//         ldp.searchBase,
//         ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//         fmt.Sprintf("(&(objectClass=posixAccount)(uid=%s))", userName),
//         []string{"dn", "uidNumber", "gidNumber", "cn", "sn", "givenName", "mail"},
//         nil,
//     )
//
//     sr, err := l.Search(searchRequest)
//     if err != nil {
//         return nil, err
//     }
//     if len(sr.Entries) != 1 {
//         return nil, fmt.Errorf("%s: %w", userName, ErrUserNotFound)
//     }
//     userDn := sr.Entries[0].DN
//
//     // rebind as the user if the password is given
//     if len(userPw) != 0 {
//         err = l.Bind(userDn, userPw)
//         if err != nil {
//             return nil, err
//         }
//     }
//
//     // create the user object
//     user := User{
//         Dn:        userDn,
//         Name:      userName,
//         UidNumber: -1,
//         GidNumber: -1,
//         Password:  "redacted",
//     }
//
//     for _, attr := range sr.Entries[0].Attributes {
//         switch attr.Name {
//         case "uidNumber":
//             uidNumber, err := strconv.Atoi(attr.Values[0])
//             if err != nil {
//                 return nil, err
//             }
//             user.UidNumber = uidNumber
//         case "gidNumber":
//             gidNumber, err := strconv.Atoi(attr.Values[0])
//             if err != nil {
//                 return nil, err
//             }
//             user.GidNumber = gidNumber
//         case "cn":
//             user.FullName = attr.Values[0]
//         case "sn":
//             user.FamilyName = attr.Values[0]
//         case "givenName":
//             user.GivenName = attr.Values[0]
//         case "mail":
//             user.Email = attr.Values[0]
//         }
//     }
//
//     // search for the groups the user is part of
//     searchRequest = ldap.NewSearchRequest(
//         ldp.searchBase,
//         ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//         fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", userName),
//         []string{"dn", "cn"},
//         nil,
//     )
//     sr, err = l.Search(searchRequest)
//     if err != nil {
//         return nil, err
//     }
//
//     if len(sr.Entries) > 0 {
//         var groups []string
//         for _, attr := range sr.Entries[0].Attributes {
//             if attr.Name == "cn" {
//                 for _, v := range attr.Values {
//                     groups = append(groups, v)
//                 }
//                 break
//             }
//         }
//         user.Groups = groups
//     }
//
//     return &user, nil
// }
//
// func (ldp *ldapClient) findGroup(groupName string) (*Group, error) {
//     l, err := ldap.Dial("tcp", ldp.server)
//     if err != nil {
//         return nil, err
//     }
//     defer l.Close()
//
//     // upgrade to tls
//     err = l.StartTLS(ldp.tlsConfig)
//     if err != nil {
//         return nil, err
//     }
//
//     // bind with the search user
//     err = l.Bind(ldp.bindDn, ldp.bindPw)
//     if err != nil {
//         return nil, err
//     }
//
//     // search for the user
//     searchRequest := ldap.NewSearchRequest(
//         ldp.searchBase,
//         ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//         fmt.Sprintf("(&(objectClass=posixGroup)(cn=%s))", groupName),
//         []string{"dn", "gidNumber"},
//         nil,
//     )
//
//     sr, err := l.Search(searchRequest)
//     if err != nil {
//         return nil, err
//     }
//     if len(sr.Entries) != 1 {
//         return nil, fmt.Errorf("%s: %w", groupName, ErrGroupNotFound)
//     }
//
//     group := &Group{
//         Dn:        sr.Entries[0].DN,
//         Name:      groupName,
//         GidNumber: -1,
//     }
//
//     for _, attr := range sr.Entries[0].Attributes {
//         if attr.Name == "gidNumber" {
//             gidNumber, err := strconv.Atoi(attr.Values[0])
//             if err != nil {
//                 return nil, err
//             }
//             group.GidNumber = gidNumber
//             break
//         }
//     }
//
//     return group, nil
// }
