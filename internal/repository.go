package course

import (
	"context"
	"fmt"
	"github.com/LuisRiveraBan/gocourse_domain/domain"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type Repository interface {
	Create(ctx context.Context, course *domain.Course) error
	GetByID(ctx context.Context, id string) (*domain.Course, error)
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, name *string, stardate, endate *time.Time) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repository struct {
	db  *gorm.DB
	log *log.Logger
}

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (r *repository) Create(ctx context.Context, course *domain.Course) error {
	if err := r.db.WithContext(ctx).Create(course).Error; err != nil {
		r.log.Println(err)
		return err
	}

	r.log.Println("user created with id: ", course.ID)
	return nil
}

func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var course []domain.Course

	tx := r.db.WithContext(ctx).Model(&course)
	tx = applyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit)
	if err := tx.Order("created_at desc").Find(&course).Error; err != nil {
		r.log.Println(err)
		return nil, err
	}
	return course, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Course, error) {

	var course domain.Course

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&course).Error; err != nil {
		r.log.Println(err)
		return nil, err
	}

	return &course, nil

}

func (r *repository) Delete(ctx context.Context, id string) error {
	course := domain.Course{ID: id}

	resul := r.db.WithContext(ctx).Delete(&course)

	if resul.Error != nil {
		r.log.Println(resul.Error)
		return resul.Error
	}

	if resul.RowsAffected == 0 {
		r.log.Printf(
			"No course found with id: %s", id,
		)

		return ErrNotFound{id}
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, name *string, stardate, endate *time.Time) error {

	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if stardate != nil {
		values["star_date"] = *stardate
	}

	if endate != nil {
		values["end_date"] = *endate
	}
	resul := r.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).Updates(values)

	if resul.Error != nil {
		r.log.Println(resul.Error)
		return resul.Error
	}

	if resul.RowsAffected == 0 {
		r.log.Printf(
			"No course found with id: %s", id,
		)
		return ErrNotFound{id}
	}
	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}
	return tx
}

func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64

	tx := r.db.WithContext(ctx).Model(&domain.Course{})

	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {

		r.log.Println(err)
		return 0, err

	}

	return int(count), nil
}
