-- name: UserInsert :one
INSERT INTO users (
    email, first_name, last_name, password, phone
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UserWhereEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UserWhereId :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: OrgInsert :one
INSERT INTO organisations (
    name, description
) VALUES ( $1, $2 )
RETURNING *;

-- name: UserAddOrg :exec
INSERT INTO user_organisations (
    user_id, org_id
) VALUES ( $1, $2 );

-- name: OrganisationWhereId :one
SELECT * FROM organisations
WHERE id = $1 LIMIT 1;

-- name: OrgWhereUser :one
SELECT org.* FROM user_organisations uo
JOIN organisations org ON uo.org_id = org.id
WHERE uo.user_id = $1 and uo.org_id = $2 limit 1;

-- name: OrgAllWhereUser :many
SELECT org.* FROM user_organisations uo
JOIN organisations org ON uo.org_id = org.id
WHERE uo.user_id = $1;

-- pretty complicated query but should work
-- gets a user if it belongs to one of another user's
-- organisation
-- name: FindUserInOrgs :one
SELECT u.* FROM users auth_user
JOIN user_organisations u_org ON u_org.user_id = auth_user.id
JOIN organisations org ON u_org.org_id = org.id
JOIN user_organisations org_users ON org_users.org_id = org.id
JOIN users u ON u.id = @find_user AND u.id = org_users.user_id
WHERE auth_user.id = @auth_user;
