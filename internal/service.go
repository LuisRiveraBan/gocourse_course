package course

import (
	"context"
	"github.com/LuisRiveraBan/gocourse_domain/domain"
	"log"
	"time"
)

type (
	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		GetByID(ctx context.Context, id string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, stardate, endate *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log: log,
		// Implement the logic to create a new user
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, name, StartDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", StartDate)

	if err != nil {
		s.log.Println("Error parsing start date:", err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)

	if err != nil {
		s.log.Println("Error parsing end date:", err)
		return nil, err
	}

	// Check if start date is after end date
	// Si starDate es mayor que endDate tira un error
	if startDateParsed.After(endDateParsed) {
		s.log.Println(ErrEndLesserStart)
		return nil, ErrEndLesserStart
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, &course); err != nil {
		s.log.Println("Error creating course:", err)
		return nil, err
	}

	return &course, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("Get All Course Service")

	course, err := s.repo.GetAll(ctx, filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.Course, error) {
	s.log.Println("Get Course By ID Service")

	course, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (s *service) Delete(ctx context.Context, id string) error {

	s.log.Println("Delete Course Service")

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id string, name *string, stardate, endate *string) error {
	s.log.Println("Update Course Service")

	var starDateParsed, endFateParsed *time.Time

	course, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if stardate != nil {
		date, err := time.Parse("2006-01-02", *stardate)

		if err != nil {
			s.log.Println(err)
			return err
		}
		// Check if start date is after end date
		// Si starDate es mayor que endDate tira un error
		if date.After(course.EndDate) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		starDateParsed = &date
	}

	if endate != nil {
		date, err := time.Parse("2006-01-02", *endate)

		if err != nil {
			s.log.Println(err)
			return err
		}
		// Check if start date is after end date
		// Si starDate es mayor que endDate tira un error
		if date.After(course.EndDate) {
			s.log.Println(ErrEndLesserStart)
			return ErrEndLesserStart
		}

		endFateParsed = &date
	}

	if err := s.repo.Update(ctx, id, name, starDateParsed, endFateParsed); err != nil {
		return err
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
