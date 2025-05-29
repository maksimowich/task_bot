package storage

import (
	"fmt"
	"sort"
	"sync"
)

type User struct {
	Id       int64
	UserName string
	ChatId   int64
}

type Task struct {
	Id          int64
	Name        string
	CreatorUser *User
	AsigneeUser *User
	IsResolved  bool
}

func (t *Task) String() string {
	var creatorUserName string
	if t.CreatorUser != nil {
		creatorUserName = t.CreatorUser.UserName
	} else {
		creatorUserName = "<no user>"
	}

	var asigneeUserName string
	if t.AsigneeUser != nil {
		asigneeUserName = t.AsigneeUser.UserName
	} else {
		asigneeUserName = "<no user>"
	}

	return fmt.Sprintf(
		"Task{Id: %d, Name: %s, CreatorUserName: %s, AsigneeUserName: %s, IsResolved: %t}",
		t.Id, t.Name, creatorUserName, asigneeUserName, t.IsResolved,
	)
}

type Storage struct {
	Users      map[int64]*User
	Tasks      map[int64]*Task
	NextTaskId int64
	Mu         sync.RWMutex
}

// ┌────────────────────────────────────────────┐
// │               User methods                 │
// └────────────────────────────────────────────┘

func (s *Storage) GetUserByIdUnsafe(
	id int64,
) (*User, error) {
	user, ok := s.Users[id]
	if !ok {
		return nil, fmt.Errorf("user does not exist")
	}

	return user, nil
}

func (s *Storage) GetUserById(
	id int64,
) (*User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.GetUserByIdUnsafe(id)
}

func (s *Storage) CreateUserUnsafe(
	id int64,
	userName string,
	chatId int64,
) (*User, error) {
	existingUser, _ := s.GetUserByIdUnsafe(id)
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	newUser := &User{
		Id:       id,
		UserName: userName,
		ChatId:   chatId,
	}
	s.Users[id] = newUser

	return newUser, nil
}

func (s *Storage) CreateUser(
	id int64,
	userName string,
	chatId int64,
) (*User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.CreateUserUnsafe(id, userName, chatId)
}

func (s *Storage) GetOrCreateUser(
	userId int64,
	userName string,
	chatId int64,
) (*User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	existingUser, _ := s.GetUserByIdUnsafe(userId)
	if existingUser != nil {
		return existingUser, nil
	}

	newUser, err := s.CreateUserUnsafe(userId, userName, chatId)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// ┌────────────────────────────────────────────┐
// │               Task methods                 │
// └────────────────────────────────────────────┘
func (s *Storage) GetTaskByIdUnsafe(
	id int64,
) (*Task, error) {
	task, ok := s.Tasks[id]
	if !ok {
		return nil, fmt.Errorf("task does not exist")
	}

	if task.IsResolved {
		return nil, fmt.Errorf("task has been resolved")
	}

	return task, nil
}

func (s *Storage) GetTaskById(
	id int64,
) (*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	return s.GetTaskByIdUnsafe(id)
}

func (s *Storage) GetTasks(
	creatorUserId int64,
	asigneeUserId int64,
) ([]*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	res := make([]*Task, 0, len(s.Tasks))
	for _, task := range s.Tasks {
		if task.IsResolved {
			continue
		}
		if creatorUserId != 0 && task.CreatorUser.Id != creatorUserId {
			continue
		}
		if asigneeUserId != 0 && task.AsigneeUser != nil && task.AsigneeUser.Id != asigneeUserId {
			continue
		}
		res = append(res, task)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Id < res[j].Id
	})
	return res, nil
}

func (s *Storage) CreateTask(
	name string,
	creatorUserId int64,
) (*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	creatorUser, err := s.GetUserByIdUnsafe(creatorUserId)
	if err != nil {
		return nil, fmt.Errorf("user lookup failed: %w", err)
	}

	taskId := s.NextTaskId
	task := &Task{
		Id:          taskId,
		Name:        name,
		CreatorUser: creatorUser,
		AsigneeUser: nil,
		IsResolved:  false,
	}

	s.Tasks[taskId] = task
	s.NextTaskId++

	return task, nil
}

func (s *Storage) AssignTask(
	id int64,
	asigneeUserId int64,
) (*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	task, err := s.GetTaskByIdUnsafe(id)
	if err != nil {
		return nil, fmt.Errorf("task lookup failed: %w", err)
	}

	asigneeUser, err := s.GetUserByIdUnsafe(asigneeUserId)
	if err != nil {
		return nil, fmt.Errorf("asignee user lookup failed: %w", err)
	}

	task.AsigneeUser = asigneeUser

	return task, nil
}

func (s *Storage) UnassignTask(
	id int64,
) (*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	task, err := s.GetTaskByIdUnsafe(id)
	if err != nil {
		return nil, fmt.Errorf("task lookup failed: %w", err)
	}

	task.AsigneeUser = nil

	return task, nil
}

func (s *Storage) ResolveTask(
	id int64,
) (*Task, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	task, err := s.GetTaskByIdUnsafe(id)
	if err != nil {
		return nil, fmt.Errorf("task lookup failed: %w", err)
	}

	task.IsResolved = true

	return task, nil
}
