// Code generated from [TypeName]HashSet.tt by T4 template. DO NOT EDIT.

//---------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool. Changes to this
//     file may cause incorrect behavior and will be lost
//     if the code is regenerated.
//
//     Generated on 2021 October 07 19:12:48 UTC
// </auto-generated>

//******************************************************************************************************
//  HashSet.go - Gbtc
//
//  Copyright © 2021, Grid Protection Alliance.	 All Rights Reserved.
//
//  Licensed to the Grid Protection Alliance (GPA) under one or more contributor license agreements. See
//  the NOTICE file distributed with this work for additional information regarding copyright ownership.
//  The GPA licenses this file to you under the MIT License (MIT), the "License"; you may not use this
//  file except in compliance with the License. You may obtain a copy of the License at:
//
//      http://opensource.org/licenses/MIT
//
//  Unless agreed to in writing, the subject software distributed under the License is distributed on an
//  "AS-IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. Refer to the
//  License for the specific language governing permissions and limitations.
//
//  Code Modification History:
//  ----------------------------------------------------------------------------------------------------
//  09/10/2021 - J. Ritchie Carroll
//       Generated original version of source code.
//
//******************************************************************************************************

package guid

type void struct{}

var member void

// TODO: Code can be changed from template to type "T" when generics are available:
//type HashSet[T any] map[T]void

// HashSet represents a distinct collection of Guid values, i.e., a set.
// A HashSet is not sorted and will not contain duplicate elements.
// The methods of the HashSet are not intrinsically thread-safe procedures,
// to guarantee thread safety, you should initiate a lock before calling a method.
type HashSet map[Guid]void

// NewHashSet creates a new set containing all elements in the specified slice.
// Returns a new HashSet that contain the specified slice items.
func NewHashSet(items []Guid) HashSet {
	hs := make(HashSet, len(items))
	hs.UnionWith(items)
	return hs
}

// Add adds the specified element to a set.
// Returns true if item was added to the set; otherwise, false.
func (hs HashSet) Add(item Guid) bool {
	if hs.Contains(item) {
		return false
	}

	hs[item] = member
	return true
}

// Remove removes the specified element from a set.
// Returns true if item was removed from the set; otherwise, false.
func (hs HashSet) Remove(item Guid) bool {
	if hs.Contains(item) {
		delete(hs, item)
		return true
	}

	return false
}

// RemoveWhere removes all elements that match the conditions defined by the specified predicate from a set.
// Returns number of elements that were removed from the set.
func (hs HashSet) RemoveWhere(predicate func(Guid) bool) int {
	var removedCount int

	for k := range hs {
		if predicate(k) && hs.Remove(k) {
			removedCount++
		}
	}

	return removedCount
}

// IsEmpty determines if set contains no elements.
// Returns true if the set is empty, i.e., has a length of zero; otherwise, false.
func (hs HashSet) IsEmpty() bool {
	return len(hs) == 0
}

// Clear removes all elements from a set.
func (hs HashSet) Clear() {
	for k := range hs {
		delete(hs, k)
	}
}

// Contains determines whether a set contains the specified element.
// Returns true if the set contains the specified element; otherwise, false.
func (hs HashSet) Contains(item Guid) bool {
	_, ok := hs[item]
	return ok
}

// Keys copies the elements of a set to a slice.
// Returns a slice containing the elements of the set.
func (hs HashSet) Keys() []Guid {
	keys := make([]Guid, len(hs))
	i := 0

	for k := range hs {
		keys[i] = k
		i++
	}

	return keys
}

// ExceptWith removes all elements in the specified slice from the current set.
func (hs HashSet) ExceptWith(other []Guid) {
	for _, v := range other {
		delete(hs, v)
	}
}

// ExceptWithSet removes all elements in the specified set from the current set.
func (hs HashSet) ExceptWithSet(other HashSet) {
	hs.ExceptWith(other.Keys())
}

// SymmetricExceptWith modifies the current set to contain only elements that are present either in the set or in the specified slice, but not both.
func (hs HashSet) SymmetricExceptWith(other []Guid) {
	if hs.IsEmpty() {
		// If set is empty, then symmetric difference is other
		hs.UnionWith(other)
		return
	}

	for _, v := range other {
		if !hs.Remove(v) {
			hs.Add(v)
		}
	}
}

// SymmetricExceptWithSet modifies the current set to contain only elements that are present either in the set or in the specified set, but not both.
func (hs HashSet) SymmetricExceptWithSet(other HashSet) {
	hs.SymmetricExceptWith(other.Keys())
}

// IntersectWith modifies the current set to contain only elements that are present in the set and in the specified slice.
func (hs HashSet) IntersectWith(other []Guid) {
	hs.IntersectWithSet(NewHashSet(other))
}

// IntersectWithSet modifies the current set to contain only elements that are present in the set and in the specified set.
func (hs HashSet) IntersectWithSet(other HashSet) {
	// Intersection of anything with empty set is empty set, so return if count is 0
	if len(hs) == 0 {
		return
	}

	if len(other) == 0 {
		hs.Clear()
		return
	}

	for k := range hs {
		if !other.Contains(k) {
			hs.Remove(k)
		}
	}
}

// UnionWith modifies the current set to contain all elements that are present in the set, the specified slice, or both.
func (hs HashSet) UnionWith(other []Guid) {
	for _, v := range other {
		hs.Add(v)
	}
}

// UnionWithSet modifies the current set to contain all elements that are present in the set, the specified set, or both.
func (hs HashSet) UnionWithSet(other HashSet) {
	hs.UnionWith(other.Keys())
}

// SetEquals determines whether a set and the specified slice contain the same elements.
// Returns true if the set is equal to other slice items; otherwise, false.
func (hs HashSet) SetEquals(other []Guid) bool {
	if len(hs) != len(other) {
		return false
	}

	for _, v := range other {
		if !hs.Contains(v) {
			return false
		}
	}

	return true
}

// SetEqualsSet determines whether a set and the specified set contain the same elements.
// Returns true if the set is equal to other set; otherwise, false.
func (hs HashSet) SetEqualsSet(other HashSet) bool {
	return hs.SetEquals(other.Keys())
}

// Overlaps determines whether the current set and a specified slice share common elements.
// Returns true if the current set and other slice items share at least one common element; otherwise, false.
func (hs HashSet) Overlaps(other []Guid) bool {
	if hs.IsEmpty() {
		return false
	}

	for _, v := range other {
		if hs.Contains(v) {
			return true
		}
	}

	return false
}

// OverlapsSet determines whether the current set and a specified set share common elements.
// Returns true if the current set and other set share at least one common element; otherwise, false.
func (hs HashSet) OverlapsSet(other HashSet) bool {
	return hs.Overlaps(other.Keys())
}

// IsSubsetOf determines whether a set is a subset of the specified slice.
// Returns true if the set is a subset of other slice items; otherwise, false.
func (hs HashSet) IsSubsetOf(other []Guid) bool {
	return hs.IsSubsetOfSet(NewHashSet(other))
}

// IsSubsetOfSet determines whether a set is a subset of the specified set.
// Returns true if the set is a subset of other set; otherwise, false.
func (hs HashSet) IsSubsetOfSet(other HashSet) bool {
	// The empty set is a subset of any set
	if hs.IsEmpty() {
		return true
	}

	// If set has more elements than slice, then it can't be a subset
	if len(hs) > len(other) {
		return false
	}

	for k := range hs {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

// IsProperSubsetOf determines whether a set is a proper subset of the specified slice.
// Returns true if the set is a proper subset of other slice items; otherwise, false.
func (hs HashSet) IsProperSubsetOf(other []Guid) bool {
	return hs.IsProperSubsetOfSet(NewHashSet(other))
}

// IsProperSubsetOfSet determines whether a set is a proper subset of the specified set.
// Returns true if the set is a proper subset of other set; otherwise, false.
func (hs HashSet) IsProperSubsetOfSet(other HashSet) bool {
	// The empty set is a proper subset of anything but the empty set
	if hs.IsEmpty() {
		return len(other) > 0
	}

	// If set has more or equal elements than slice, then it can't be a proper subset
	if len(hs) >= len(other) {
		return false
	}

	for k := range hs {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

// IsSupersetOf determines whether a set is a superset of the specified slice.
// Returns true if the set is a superset of other slice items; otherwise, false.
func (hs HashSet) IsSupersetOf(other []Guid) bool {
	// If other is the empty set then this is a superset
	if len(other) == 0 {
		return true
	}

	// If slice has more elements than set, then it can't be a superset
	if len(other) > len(hs) {
		return false
	}

	for _, v := range other {
		if !hs.Contains(v) {
			return false
		}
	}

	return true
}

// IsSupersetOfSet determines whether a set is a superset of the specified set.
// Returns true if the set is a superset of other set; otherwise, false.
func (hs HashSet) IsSupersetOfSet(other HashSet) bool {
	return hs.IsSupersetOf(other.Keys())
}

// IsProperSupersetOf determines whether a set is a proper superset of the specified slice.
// Returns true if the set is a proper superset of other slice items; otherwise, false.
func (hs HashSet) IsProperSupersetOf(other []Guid) bool {
	// The empty set is not a proper subset of any set
	if hs.IsEmpty() {
		return true
	}

	// If other is the empty set then this is a proper superset
	if len(other) == 0 {
		return true // Set has at least one element, based on prior check
	}

	// If slice has more or equal elements than set, then it can't be a proper superset
	if len(other) >= len(hs) {
		return false
	}

	for _, v := range other {
		if !hs.Contains(v) {
			return false
		}
	}

	return true
}

// IsProperSupersetOfSet determines whether a set is a proper superset of the specified set.
// Returns true if the set is a proper superset of other set; otherwise, false.
func (hs HashSet) IsProperSupersetOfSet(other HashSet) bool {
	return hs.IsProperSupersetOf(other.Keys())
}
