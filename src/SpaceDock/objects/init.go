/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import "SpaceDock"

/*
 This function creates tables for all datatypes
 */
func init() {
    SpaceDock.CreateTable(&User{})
    SpaceDock.CreateTable(&Role{})
    SpaceDock.CreateTable(&RoleUser{})
}