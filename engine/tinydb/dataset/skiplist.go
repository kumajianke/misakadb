package dataset

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	defaultMaxLevel    = 16
	defaultProbability = 0.5
)

type PassStuct struct {
	Name string
}

type SkipListNode[T any] struct {
	Key     int
	Data    []T
	forward []*SkipListNode[T]
}

type SkipList[T any] struct {
	head        *SkipListNode[T]
	level       int
	length      int
	maxLevel    int
	probability float64
	rng         *rand.Rand
}

type KVPair[T any] struct {
	Key  int
	Data []T
}

func NewSkipList[T any]() *SkipList[T] {
	return NewSkipListWithConfig[T](defaultMaxLevel, defaultProbability)
}

func NewSkipListWithConfig[T any](maxLevel int, probability float64) *SkipList[T] {
	if maxLevel <= 0 {
		maxLevel = defaultMaxLevel
	}
	if probability <= 0 || probability >= 1 {
		probability = defaultProbability
	}

	return &SkipList[T]{
		head: &SkipListNode[T]{
			Key:     -1,
			forward: make([]*SkipListNode[T], maxLevel),
		},
		level:       1,
		maxLevel:    maxLevel,
		probability: probability,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewPassSkipList() *SkipList[PassStuct] {
	return NewSkipList[PassStuct]()
}

func (skipList *SkipList[T]) Len() int {
	return skipList.length
}

func (skipList *SkipList[T]) randomLevel() int {
	level := 1
	for level < skipList.maxLevel && skipList.rng.Float64() < skipList.probability {
		level++
	}
	return level
}

func (skipList *SkipList[T]) findUpdatePath(key int) ([]*SkipListNode[T], *SkipListNode[T]) {
	update := make([]*SkipListNode[T], skipList.maxLevel)
	current := skipList.head

	for index := skipList.level - 1; index >= 0; index-- {
		for current.forward[index] != nil && current.forward[index].Key < key {
			current = current.forward[index]
		}
		update[index] = current
	}

	candidate := current.forward[0]
	return update, candidate
}

func (skipList *SkipList[T]) Insert(key int, data T) {
	update, candidate := skipList.findUpdatePath(key)
	if candidate != nil && candidate.Key == key {
		candidate.Data = append(candidate.Data, data)
		skipList.length++
		return
	}

	nodeLevel := skipList.randomLevel()
	if nodeLevel > skipList.level {
		for index := skipList.level; index < nodeLevel; index++ {
			update[index] = skipList.head
		}
		skipList.level = nodeLevel
	}

	newNode := &SkipListNode[T]{
		Key:     key,
		Data:    []T{data},
		forward: make([]*SkipListNode[T], nodeLevel),
	}

	for index := range nodeLevel {
		newNode.forward[index] = update[index].forward[index]
		update[index].forward[index] = newNode
	}

	skipList.length++
}

func (skipList *SkipList[T]) Get(key int) ([]T, bool) {
	current := skipList.head

	for index := skipList.level - 1; index >= 0; index-- {
		for current.forward[index] != nil && current.forward[index].Key < key {
			current = current.forward[index]
		}
	}

	current = current.forward[0]
	if current != nil && current.Key == key {
		return current.Data, true
	}

	return nil, false
}

func (skipList *SkipList[T]) Delete(key int) bool {
	update, candidate := skipList.findUpdatePath(key)
	if candidate == nil || candidate.Key != key {
		return false
	}

	for index := range skipList.level {
		if update[index].forward[index] != candidate {
			continue
		}
		update[index].forward[index] = candidate.forward[index]
	}

	for skipList.level > 1 && skipList.head.forward[skipList.level-1] == nil {
		skipList.level--
	}

	skipList.length -= len(candidate.Data)
	return true
}

func (skipList *SkipList[T]) Contains(key int) bool {
	values, found := skipList.Get(key)
	return found && len(values) > 0
}

func (skipList *SkipList[T]) Items() []KVPair[T] {
	items := make([]KVPair[T], 0, skipList.length)
	current := skipList.head.forward[0]

	for current != nil {
		items = append(items, KVPair[T]{
			Key:  current.Key,
			Data: current.Data,
		})
		current = current.forward[0]
	}

	return items
}

func (skipList *SkipList[T]) GetWith(operator string, value any, limit int) ([]KVPair[T], error) {
	return skipList.GetWithPage(operator, value, 0, limit)
}

func (skipList *SkipList[T]) GetWithPage(operator string, value any, offset int, limit int) ([]KVPair[T], error) {
	normalized := strings.ToLower(strings.TrimSpace(operator))
	switch normalized {
	case "<", "<=":
		target, err := toInt(value)
		if err != nil {
			return nil, err
		}
		return skipList.collectUntilKey(target, normalized == "<=", offset, limit), nil
	case ">", ">=", "=", "==":
		target, err := toInt(value)
		if err != nil {
			return nil, err
		}
		return skipList.collectFromKey(target, normalized, offset, limit), nil
	case "like":
		pattern := fmt.Sprint(value)
		return skipList.filterByKey(func(key int) bool {
			return likeMatch(fmt.Sprintf("%d", key), pattern)
		}, offset, limit), nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func (skipList *SkipList[T]) DeleteWith(operator string, value any, limit int) (int, error) {
	normalized := strings.ToLower(strings.TrimSpace(operator))
	switch normalized {
	case ">", ">=", "=", "==":
		target, err := toInt(value)
		if err != nil {
			return 0, err
		}
		return skipList.deleteFromKey(target, normalized, limit), nil
	}

	rows, err := skipList.GetWith(operator, value, limit)
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, row := range rows {
		if skipList.Delete(row.Key) {
			deleted++
		}
	}

	return deleted, nil
}

func (skipList *SkipList[T]) deleteFromKey(target int, operator string, limit int) int {
	start := skipList.findFirstGreaterOrEqual(target)
	if start == nil {
		return 0
	}

	switch operator {
	case ">":
		for start != nil && start.Key <= target {
			start = start.forward[0]
		}
	case "=", "==":
		if start.Key != target {
			return 0
		}
	}

	if start == nil {
		return 0
	}

	update, candidate := skipList.findUpdatePath(start.Key)
	if candidate != start {
		return 0
	}

	deleted := 0
	current := start
	for current != nil {
		if !matchForwardOperator(current.Key, target, operator) {
			break
		}

		next := current.forward[0]
		for index := range skipList.level {
			if update[index].forward[index] != current {
				continue
			}
			update[index].forward[index] = current.forward[index]
		}

		skipList.length -= len(current.Data)
		deleted++
		if limit > 0 && deleted >= limit {
			break
		}
		current = next
	}

	for skipList.level > 1 && skipList.head.forward[skipList.level-1] == nil {
		skipList.level--
	}

	return deleted
}

func (skipList *SkipList[T]) collectFromKey(target int, operator string, offset int, limit int) []KVPair[T] {
	start := skipList.findFirstGreaterOrEqual(target)
	if start == nil {
		return nil
	}

	switch operator {
	case ">":
		for start != nil && start.Key <= target {
			start = start.forward[0]
		}
	case "=", "==":
		if start.Key != target {
			return nil
		}
	}

	return collectForward(start, func(key int) bool {
		switch operator {
		case ">", ">=":
			return true
		case "=", "==":
			return key == target
		default:
			return false
		}
	}, offset, limit, true)
}

func (skipList *SkipList[T]) collectUntilKey(target int, includeEqual bool, offset int, limit int) []KVPair[T] {
	return collectForward(skipList.head.forward[0], func(key int) bool {
		if includeEqual {
			return key <= target
		}
		return key < target
	}, offset, limit, true)
}

func (skipList *SkipList[T]) findFirstGreaterOrEqual(key int) *SkipListNode[T] {
	current := skipList.head
	for index := skipList.level - 1; index >= 0; index-- {
		for current.forward[index] != nil && current.forward[index].Key < key {
			current = current.forward[index]
		}
	}
	return current.forward[0]
}

func (skipList *SkipList[T]) filterByKey(match func(int) bool, offset int, limit int) []KVPair[T] {
	return collectForward(skipList.head.forward[0], match, offset, limit, false)
}

func matchForwardOperator(key int, target int, operator string) bool {
	switch operator {
	case ">", ">=":
		return true
	case "=", "==":
		return key == target
	default:
		return false
	}
}

func collectForward[T any](start *SkipListNode[T], match func(int) bool, offset int, limit int, stopOnFirstMiss bool) []KVPair[T] {
	items := make([]KVPair[T], 0)
	current := start
	if offset < 0 {
		offset = 0
	}
	skipped := 0

	for current != nil {
		if match(current.Key) {
			if skipped < offset {
				skipped++
				current = current.forward[0]
				continue
			}
			items = append(items, KVPair[T]{
				Key:  current.Key,
				Data: current.Data,
			})
			if limit > 0 && len(items) >= limit {
				break
			}
		} else if stopOnFirstMiss && (len(items) > 0 || skipped > 0) {
			break
		}
		current = current.forward[0]
	}

	return items
}

func toInt(value any) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int8:
		return int(typed), nil
	case int16:
		return int(typed), nil
	case int32:
		return int(typed), nil
	case int64:
		return int(typed), nil
	case uint:
		return int(typed), nil
	case uint8:
		return int(typed), nil
	case uint16:
		return int(typed), nil
	case uint32:
		return int(typed), nil
	case uint64:
		return int(typed), nil
	case string:
		var parsed int
		_, err := fmt.Sscanf(strings.TrimSpace(typed), "%d", &parsed)
		if err != nil {
			return 0, fmt.Errorf("value %q can not convert to int", typed)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("unsupported key type: %T", value)
	}
}

func likeMatch(value string, pattern string) bool {
	return likeMatchRunes([]rune(value), []rune(pattern))
}

func likeMatchRunes(value []rune, pattern []rune) bool {
	if len(pattern) == 0 {
		return len(value) == 0
	}

	if pattern[0] == '%' {
		if likeMatchRunes(value, pattern[1:]) {
			return true
		}
		if len(value) > 0 {
			return likeMatchRunes(value[1:], pattern)
		}
		return false
	}

	if len(value) == 0 {
		return false
	}

	if value[0] != pattern[0] {
		return false
	}

	return likeMatchRunes(value[1:], pattern[1:])
}
