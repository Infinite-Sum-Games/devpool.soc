-- name: GetInstallationIdQuery :one

SELECT installation_id
FROM repository
WHERE url = $1 ;
