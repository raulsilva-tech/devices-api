package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raulsilva-tech/devices-api/internal/domain"
	"github.com/stretchr/testify/require"
)

// --- Mock repository (manual, lightweight) ---
type mockDeviceRepo struct {
	CreateDeviceFunc      func(ctx context.Context, device *domain.Device) (string, error)
	UpdateDeviceFunc      func(ctx context.Context, device *domain.Device) error
	DeleteDeviceFunc      func(ctx context.Context, id string) error
	GetDeviceByIdFunc     func(ctx context.Context, id string) (*domain.Device, error)
	GetDevicesFunc        func(ctx context.Context) ([]domain.Device, error)
	GetDevicesByBrandFunc func(ctx context.Context, brand string) ([]domain.Device, error)
	GetDevicesByStateFunc func(ctx context.Context, state string) ([]domain.Device, error)
}

func (m *mockDeviceRepo) CreateDevice(ctx context.Context, device *domain.Device) (string, error) {
	return m.CreateDeviceFunc(ctx, device)
}
func (m *mockDeviceRepo) UpdateDevice(ctx context.Context, device *domain.Device) error {
	return m.UpdateDeviceFunc(ctx, device)
}
func (m *mockDeviceRepo) DeleteDevice(ctx context.Context, id string) error {
	return m.DeleteDeviceFunc(ctx, id)
}
func (m *mockDeviceRepo) GetDeviceById(ctx context.Context, id string) (*domain.Device, error) {
	return m.GetDeviceByIdFunc(ctx, id)
}
func (m *mockDeviceRepo) GetDevices(ctx context.Context) ([]domain.Device, error) {
	return m.GetDevicesFunc(ctx)
}
func (m *mockDeviceRepo) GetDevicesByBrand(ctx context.Context, brand string) ([]domain.Device, error) {
	return m.GetDevicesByBrandFunc(ctx, brand)
}
func (m *mockDeviceRepo) GetDevicesByState(ctx context.Context, state string) ([]domain.Device, error) {
	return m.GetDevicesByStateFunc(ctx, state)
}

// --- helpers ---
func makeDeviceWithState(state domain.DeviceState) *domain.Device {
	id := uuid.New().String()
	d, _ := domain.NewDevice(id, "Name", "Brand", state, time.Now())
	return d
}

func deviceServiceWithMock(m *mockDeviceRepo) *DeviceService {
	return NewDeviceService(m)
}

// -------------------- Tests --------------------

func TestCreateDevice_SuccessAndRepoError(t *testing.T) {
	ctx := context.Background()

	// success case
	mock := &mockDeviceRepo{
		CreateDeviceFunc: func(ctx context.Context, device *domain.Device) (string, error) {
			// emulate repo returning the same id
			return device.ID, nil
		},
	}
	svc := deviceServiceWithMock(mock)

	id, err := svc.CreateDevice(ctx, CreateDeviceInput{
		Name:  "Device A",
		Brand: "Brand A",
		State: domain.DeviceAvailable,
	})
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// repo error case
	mockErr := errors.New("db error")
	mock2 := &mockDeviceRepo{
		CreateDeviceFunc: func(ctx context.Context, device *domain.Device) (string, error) {
			return "", mockErr
		},
	}
	svc2 := deviceServiceWithMock(mock2)
	_, err = svc2.CreateDevice(ctx, CreateDeviceInput{
		Name:  "Device A",
		Brand: "Brand A",
		State: domain.DeviceAvailable,
	})
	require.ErrorIs(t, err, mockErr)
}

func TestUpdateDevice_InUse_AllowsOnlyName(t *testing.T) {
	ctx := context.Background()

	orig := makeDeviceWithState(domain.DeviceInUse)
	orig.Name = "Old"
	orig.Brand = "OrigBrand"

	var updatedSaved *domain.Device

	mock := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return orig, nil
		},
		UpdateDeviceFunc: func(ctx context.Context, device *domain.Device) error {
			// capture what was saved
			copy := *device
			updatedSaved = &copy
			return nil
		},
	}

	svc := deviceServiceWithMock(mock)

	out, err := svc.UpdateDevice(ctx, UpdateDeviceInput{
		ID:    orig.ID,
		Name:  "NewName",
		Brand: "NewBrand",             // should be ignored
		State: domain.DeviceAvailable, // should be ignored
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	// only name changed
	require.Equal(t, "NewName", out.Device.Name)
	require.Equal(t, "OrigBrand", out.Device.Brand)
	require.Equal(t, domain.DeviceInUse, out.Device.State)

	// persisted values
	require.NotNil(t, updatedSaved)
	require.Equal(t, "NewName", updatedSaved.Name)
	require.Equal(t, "OrigBrand", updatedSaved.Brand)
	require.Equal(t, domain.DeviceInUse, updatedSaved.State)

	// ensure ignored fields list contains brand and state
	require.Contains(t, out.IgnoredFields, "brand")
	require.Contains(t, out.IgnoredFields, "state")
	require.Contains(t, out.UpdatedFields, "name")
}

func TestUpdateDevice_NotInUse_AllFieldsChange(t *testing.T) {
	ctx := context.Background()

	orig := makeDeviceWithState(domain.DeviceAvailable)
	orig.Name = "Old"
	orig.Brand = "OrigBrand"

	var updatedSaved *domain.Device

	mock := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return orig, nil
		},
		UpdateDeviceFunc: func(ctx context.Context, device *domain.Device) error {
			copy := *device
			updatedSaved = &copy
			return nil
		},
	}

	svc := deviceServiceWithMock(mock)

	out, err := svc.UpdateDevice(ctx, UpdateDeviceInput{
		ID:    orig.ID,
		Name:  "NewName",
		Brand: "NewBrand",
		State: domain.DeviceInUse,
	})
	require.NoError(t, err)
	require.NotNil(t, out)

	require.Equal(t, "NewName", out.Device.Name)
	require.Equal(t, "NewBrand", out.Device.Brand)
	require.Equal(t, domain.DeviceInUse, out.Device.State)

	// persisted equals output
	require.NotNil(t, updatedSaved)
	require.Equal(t, out.Device.Name, updatedSaved.Name)
	require.Equal(t, out.Device.Brand, updatedSaved.Brand)
	require.Equal(t, out.Device.State, updatedSaved.State)
}

func TestUpdateDevice_InvalidStateProvided(t *testing.T) {
	ctx := context.Background()

	orig := makeDeviceWithState(domain.DeviceAvailable)

	mock := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return orig, nil
		},
	}

	svc := deviceServiceWithMock(mock)

	_, err := svc.UpdateDevice(ctx, UpdateDeviceInput{
		ID:    orig.ID,
		Name:  "X",
		Brand: "B",
		State: domain.DeviceState("invalid-state"),
	})
	require.ErrorIs(t, err, domain.ErrInvalidState)
}

func TestUpdateDevice_RepoGetErrorAndRepoUpdateError(t *testing.T) {
	ctx := context.Background()

	getErr := errors.New("get failed")
	mockGetErr := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return nil, getErr
		},
	}
	svcGetErr := deviceServiceWithMock(mockGetErr)
	_, err := svcGetErr.UpdateDevice(ctx, UpdateDeviceInput{ID: "x"})
	require.ErrorIs(t, err, getErr)

	// update error
	orig := makeDeviceWithState(domain.DeviceAvailable)
	mockUpdateErr := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return orig, nil
		},
		UpdateDeviceFunc: func(ctx context.Context, device *domain.Device) error {
			return errors.New("update failed")
		},
	}
	svcUpdateErr := deviceServiceWithMock(mockUpdateErr)
	_, err = svcUpdateErr.UpdateDevice(ctx, UpdateDeviceInput{
		ID:    orig.ID,
		Name:  "New",
		Brand: "B",
		State: domain.DeviceInUse,
	})
	require.Error(t, err)
}

func TestGetDeviceById_SuccessAndNotFound(t *testing.T) {
	ctx := context.Background()

	orig := makeDeviceWithState(domain.DeviceAvailable)
	mock := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return orig, nil
		},
	}
	svc := deviceServiceWithMock(mock)

	out, err := svc.GetDeviceById(ctx, orig.ID)
	require.NoError(t, err)
	require.Equal(t, orig.ID, out.ID)

	// not found
	mockNotFound := &mockDeviceRepo{
		GetDeviceByIdFunc: func(ctx context.Context, id string) (*domain.Device, error) {
			return nil, sql.ErrNoRows
		},
	}
	svcNotFound := deviceServiceWithMock(mockNotFound)
	_, err = svcNotFound.GetDeviceById(ctx, "x")
	require.Error(t, err)
	require.Contains(t, err.Error(), "device id x not found")
}

func TestGetDevices_ListAndEmpty(t *testing.T) {
	ctx := context.Background()

	dev1 := *makeDeviceWithState(domain.DeviceAvailable)
	dev1.Name = "A"
	dev2 := *makeDeviceWithState(domain.DeviceInactive)
	dev2.Name = "B"

	mock := &mockDeviceRepo{
		GetDevicesFunc: func(ctx context.Context) ([]domain.Device, error) {
			return []domain.Device{dev1, dev2}, nil
		},
	}
	svc := deviceServiceWithMock(mock)

	list, err := svc.GetDevices(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, "A", list[0].Name)
	require.Equal(t, "B", list[1].Name)

	// empty list
	mockEmpty := &mockDeviceRepo{
		GetDevicesFunc: func(ctx context.Context) ([]domain.Device, error) {
			return []domain.Device{}, nil
		},
	}
	svcEmpty := deviceServiceWithMock(mockEmpty)
	list2, err := svcEmpty.GetDevices(ctx)
	require.NoError(t, err)
	require.Len(t, list2, 0)
}

func TestGetDevicesByBrandAndState(t *testing.T) {
	ctx := context.Background()

	dev := *makeDeviceWithState(domain.DeviceAvailable)
	dev.Brand = "Acme"

	mockByBrand := &mockDeviceRepo{
		GetDevicesByBrandFunc: func(ctx context.Context, brand string) ([]domain.Device, error) {
			require.Equal(t, "Acme", brand)
			return []domain.Device{dev}, nil
		},
	}
	svcBrand := deviceServiceWithMock(mockByBrand)
	list, err := svcBrand.GetDevicesByBrand(ctx, "Acme")
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Acme", list[0].Brand)

	mockByState := &mockDeviceRepo{
		GetDevicesByStateFunc: func(ctx context.Context, state string) ([]domain.Device, error) {
			require.Equal(t, string(domain.DeviceAvailable), state)
			return []domain.Device{dev}, nil
		},
	}
	svcState := deviceServiceWithMock(mockByState)
	list2, err := svcState.GetDevicesByState(ctx, string(domain.DeviceAvailable))
	require.NoError(t, err)
	require.Len(t, list2, 1)
	require.Equal(t, domain.DeviceAvailable, list2[0].State)
}
