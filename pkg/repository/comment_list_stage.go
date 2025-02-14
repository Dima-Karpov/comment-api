package repository

import (
	"bufio"
	"comment-api/internal/domain"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LocalFile struct {
	Path string
}

type CommentListStage struct {
	stage *LocalFile
}

func NewCommentListStage(stage *LocalFile) *CommentListStage {
	return &CommentListStage{stage: stage}
}

func (s *CommentListStage) Create(input domain.CommentList) (uuid.UUID, error) {
	var comment domain.Comment
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.Description = input.Description

	if input.ParentID != "" {
		comment.ParentID = input.ParentID
		err := saveCommentToFile(s.stage.Path, comment)
		if err != nil {
			return uuid.Nil, err
		}

		return comment.ID, nil
	}
	comment.ID = uuid.New()
	err := saveCommentToFile(s.stage.Path, comment)
	if err != nil {
		return uuid.Nil, err
	}

	return comment.ID, nil
}

func (s *CommentListStage) GetById(id uuid.UUID) (domain.Comment, error) {
	var comment domain.Comment

	file, err := os.Open(s.stage.Path)
	if err != nil {
		return comment, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var c domain.Comment
		err := json.Unmarshal(scanner.Bytes(), &c)
		if err != nil {
			return comment, fmt.Errorf("failed to unmarshal file: %w", err)
		}

		if c.ID == id {
			return c, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return comment, fmt.Errorf("failed to scan file: %w", err)
	}

	return comment, fmt.Errorf("comment with id %s not found", id)
}

func (s *CommentListStage) Delete(id uuid.UUID) error {
	file, err := os.Open(s.stage.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var comments []domain.Comment
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var c domain.Comment
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return fmt.Errorf("failed to unmarshal file: %w", err)
		}

		if c.ID != id { // Оставляем все комментарии, КРОМЕ того, который надо удалить
			comments = append(comments, c)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	// Перезаписываем файл без удаленного комментария
	tempFilePath := s.stage.Path + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	for _, c := range comments {
		commentJSON, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal comment: %w", err)
		}
		if _, err := tempFile.Write(append(commentJSON, '\n')); err != nil {
			return fmt.Errorf("failed to write to temp file: %w", err)
		}
	}

	// Перемещаем временный файл на место старого
	if err := os.Rename(tempFilePath, s.stage.Path); err != nil {
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}

func (s *CommentListStage) Update(id uuid.UUID, list domain.UpdateCommentList) error {
	file, err := os.Open(s.stage.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var comments []domain.Comment
	scanner := bufio.NewScanner(file)
	found := false

	for scanner.Scan() {
		var c domain.Comment
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return fmt.Errorf("failed to unmarshal comment: %w", err)
		}

		if c.ID == id {
			// Обновляем только Description and UpdateAt
			c.Description = list.Description
			c.UpdatedAt = time.Now()
			found = true
		}

		comments = append(comments, c)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}
	if !found {
		return fmt.Errorf("comment with id %s not found", id)
	}

	// Перезаписываем файл с обновленным комментарием
	tempFilePath := s.stage.Path + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	for _, c := range comments {
		commentJSON, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal comment: %w", err)
		}
		if _, err := tempFile.Write(append(commentJSON, '\n')); err != nil {
			return fmt.Errorf("failed to write to temp file: %w", err)
		}
	}

	// Перемещаем временный файд на место старого
	if err := os.Rename(tempFilePath, s.stage.Path); err != nil {
		return fmt.Errorf("failed to replace original file: %w", err)
	}
	return nil
}

func saveCommentToFile(filePath string, comment domain.Comment) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("failed to open file: %s", err)
		return err
	}
	defer file.Close()

	commentJSON, err := json.Marshal(comment)
	if err != nil {
		logrus.Fatalf("failed to marshal comment to JSON: %s", err)
		return err
	}
	if _, err := file.Write(append(commentJSON, '\n')); err != nil {
		logrus.Fatalf("failed to write to file: %s", err)
		return err
	}
	fmt.Println("Comment saved successfully.")

	return nil
}
