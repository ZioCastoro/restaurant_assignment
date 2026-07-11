package postgre_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ZioCastoro/restaurant_assignment/postgre"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Connect(t *testing.T) {
	t.Parallel()

	// initialize input arguments.
	type args struct {
		ctx    context.Context
		driver string
		dsn    string
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T) (args, testDB)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "happy flow",
			setup: func(t *testing.T) (args, testDB) {
				dsn, mockDB := setDriverMockConnect(t, func(mock sqlmock.Sqlmock) {
					mock.ExpectPing()
				})

				return args{
					ctx:    context.Background(),
					driver: "sqlmock",
					dsn:    dsn,
				}, mockDB
			},
			wantErr: false,
		},
		{
			name: "connection error",
			setup: func(t *testing.T) (args, testDB) {
				return args{
					ctx:    context.Background(),
					driver: "invalid_driver",
					dsn:    "invalid_dsn",
				}, testDB{}
			},
			wantErr: true,
			errMsg:  "sql: unknown driver",
		},
		{
			name: "ping error",
			setup: func(t *testing.T) (args, testDB) {
				dsn, mockDB := setDriverMockConnect(t, func(mock sqlmock.Sqlmock) {
					mock.ExpectPing().WillReturnError(errors.New("ping error"))
				})

				return args{
					ctx:    context.Background(),
					driver: "sqlmock",
					dsn:    dsn,
				}, mockDB
			},
			wantErr: true,
			errMsg:  "db.PingContext(): ping error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testArgs, mockDB := tt.setup(t)

			t.Cleanup(func() {
				if mockDB.DB != nil {
					_ = mockDB.Close()
				}
			})

			var err error
			// call method under test with panic recover.
			got, err := func(args args) (*sql.DB, error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered panic: %v", r)
					}
				}()

				return postgre.Connect(args.ctx, args.driver, args.dsn)
			}(testArgs)

			t.Cleanup(func() {
				if got != nil {
					_ = got.Close()
				}
			})
			
			// assert error.
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			// assert result.
			if !tt.wantErr {
				require.NotNil(t, got)
			} else {
				require.Nil(t, got)
			}

			// assert mock expectations.
			if mockDB.mock != nil {
				assert.NoError(t, mockDB.mock.ExpectationsWereMet())
			}
		})
	}
}

func setDriverMockConnect(t *testing.T, setup func(sqlmock.Sqlmock)) (string, testDB) {
	t.Helper()

	dsn := t.Name()
	db, mock, err := sqlmock.NewWithDSN(dsn, sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	setup(mock)

	return dsn, testDB{DB: db, mock: mock}
}
