package database

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/canonical/microcluster/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var serviceObjects = cluster.RegisterStmt(`
SELECT services.id, internal_cluster_members.name AS member, services.service
  FROM services
  JOIN internal_cluster_members ON services.member_id = internal_cluster_members.id
  ORDER BY internal_cluster_members.id, services.service
`)

var serviceObjectsByMember = cluster.RegisterStmt(`
SELECT services.id, internal_cluster_members.name AS member, services.service
  FROM services
  JOIN internal_cluster_members ON services.member_id = internal_cluster_members.id
  WHERE ( member = ? )
  ORDER BY internal_cluster_members.id, services.service
`)

var serviceObjectsByService = cluster.RegisterStmt(`
SELECT services.id, internal_cluster_members.name AS member, services.service
  FROM services
  JOIN internal_cluster_members ON services.member_id = internal_cluster_members.id
  WHERE ( services.service = ? )
  ORDER BY internal_cluster_members.id, services.service
`)

var serviceObjectsByMemberAndService = cluster.RegisterStmt(`
SELECT services.id, internal_cluster_members.name AS member, services.service
  FROM services
  JOIN internal_cluster_members ON services.member_id = internal_cluster_members.id
  WHERE ( member = ? AND services.service = ? )
  ORDER BY internal_cluster_members.id, services.service
`)

var serviceID = cluster.RegisterStmt(`
SELECT services.id FROM services
  JOIN internal_cluster_members ON services.member_id = internal_cluster_members.id
  WHERE internal_cluster_members.name = ? AND services.service = ?
`)

var serviceCreate = cluster.RegisterStmt(`
INSERT INTO services (member_id, service)
  VALUES ((SELECT internal_cluster_members.id FROM internal_cluster_members WHERE internal_cluster_members.name = ?), ?)
`)

var serviceDeleteByMember = cluster.RegisterStmt(`
DELETE FROM services WHERE member_id = (SELECT internal_cluster_members.id FROM internal_cluster_members WHERE internal_cluster_members.name = ?)
`)

var serviceDeleteByMemberAndService = cluster.RegisterStmt(`
DELETE FROM services WHERE member_id = (SELECT internal_cluster_members.id FROM internal_cluster_members WHERE internal_cluster_members.name = ?) AND service = ?
`)

var serviceUpdate = cluster.RegisterStmt(`
UPDATE services
  SET member_id = (SELECT internal_cluster_members.id FROM internal_cluster_members WHERE internal_cluster_members.name = ?), service = ?
 WHERE id = ?
`)

// serviceColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Service entity.
func serviceColumns() string {
	return "services.id, internal_cluster_members.name AS member, services.service"
}

// getServices can be used to run handwritten sql.Stmts to return a slice of objects.
func getServices(ctx context.Context, stmt *sql.Stmt, args ...any) ([]Service, error) {
	objects := make([]Service, 0)

	dest := func(scan func(dest ...any) error) error {
		s := Service{}
		err := scan(&s.ID, &s.Member, &s.Service)
		if err != nil {
			return err
		}

		objects = append(objects, s)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"services\" table: %w", err)
	}

	return objects, nil
}

// getServices can be used to run handwritten query strings to return a slice of objects.
func getServicesRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]Service, error) {
	objects := make([]Service, 0)

	dest := func(scan func(dest ...any) error) error {
		s := Service{}
		err := scan(&s.ID, &s.Member, &s.Service)
		if err != nil {
			return err
		}

		objects = append(objects, s)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"services\" table: %w", err)
	}

	return objects, nil
}

// GetServices returns all available services.
// generator: service GetMany
func GetServices(ctx context.Context, tx *sql.Tx, filters ...ServiceFilter) ([]Service, error) {
	var err error

	// Result slice.
	objects := make([]Service, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = cluster.Stmt(tx, serviceObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"serviceObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Member != nil && filter.Service != nil {
			args = append(args, []any{filter.Member, filter.Service}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, serviceObjectsByMemberAndService)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"serviceObjectsByMemberAndService\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(serviceObjectsByMemberAndService)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"serviceObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Service != nil && filter.Member == nil {
			args = append(args, []any{filter.Service}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, serviceObjectsByService)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"serviceObjectsByService\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(serviceObjectsByService)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"serviceObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Member != nil && filter.Service == nil {
			args = append(args, []any{filter.Member}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, serviceObjectsByMember)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"serviceObjectsByMember\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(serviceObjectsByMember)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"serviceObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Member == nil && filter.Service == nil {
			return nil, fmt.Errorf("Cannot filter on empty ServiceFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getServices(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getServicesRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"services\" table: %w", err)
	}

	return objects, nil
}

// GetService returns the service with the given key.
// generator: service GetOne
func GetService(ctx context.Context, tx *sql.Tx, member string, service string) (*Service, error) {
	filter := ServiceFilter{}
	filter.Member = &member
	filter.Service = &service

	objects, err := GetServices(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"services\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "Service not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"services\" entry matches")
	}
}

// GetServiceID return the ID of the service with the given key.
// generator: service ID
func GetServiceID(ctx context.Context, tx *sql.Tx, member string, service string) (int64, error) {
	stmt, err := cluster.Stmt(tx, serviceID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"serviceID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, member, service)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "Service not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"services\" ID: %w", err)
	}

	return id, nil
}

// ServiceExists checks if a service with the given key exists.
// generator: service Exists
func ServiceExists(ctx context.Context, tx *sql.Tx, member string, service string) (bool, error) {
	_, err := GetServiceID(ctx, tx, member, service)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateService adds a new service to the database.
// generator: service Create
func CreateService(ctx context.Context, tx *sql.Tx, object Service) (int64, error) {
	// Check if a service with the same key exists.
	exists, err := ServiceExists(ctx, tx, object.Member, object.Service)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"services\" entry already exists")
	}

	args := make([]any, 2)

	// Populate the statement arguments.
	args[0] = object.Member
	args[1] = object.Service

	// Prepared statement to use.
	stmt, err := cluster.Stmt(tx, serviceCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"serviceCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"services\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"services\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteService deletes the service matching the given key parameters.
// generator: service DeleteOne-by-Member-and-Service
func DeleteService(ctx context.Context, tx *sql.Tx, member string, service string) error {
	stmt, err := cluster.Stmt(tx, serviceDeleteByMemberAndService)
	if err != nil {
		return fmt.Errorf("Failed to get \"serviceDeleteByMemberAndService\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(member, service)
	if err != nil {
		return fmt.Errorf("Delete \"services\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "Service not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d Service rows instead of 1", n)
	}

	return nil
}

// DeleteServices deletes the service matching the given key parameters.
// generator: service DeleteMany-by-Member
func DeleteServices(ctx context.Context, tx *sql.Tx, member string) error {
	stmt, err := cluster.Stmt(tx, serviceDeleteByMember)
	if err != nil {
		return fmt.Errorf("Failed to get \"serviceDeleteByMember\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(member)
	if err != nil {
		return fmt.Errorf("Delete \"services\": %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}

// UpdateService updates the service matching the given key parameters.
// generator: service Update
func UpdateService(ctx context.Context, tx *sql.Tx, member string, service string, object Service) error {
	id, err := GetServiceID(ctx, tx, member, service)
	if err != nil {
		return err
	}

	stmt, err := cluster.Stmt(tx, serviceUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"serviceUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Member, object.Service, id)
	if err != nil {
		return fmt.Errorf("Update \"services\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}