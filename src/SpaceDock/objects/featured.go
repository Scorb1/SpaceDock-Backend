/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
)

type Featured struct {
    Model

    Mod    Mod `json:"mod" spacedock:"lock"`
    ModID  uint `json:"mod_id" spacedock:"lock"`
}

func (s *Featured) AfterFind() {
    if SpaceDock.DBRecursion == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.DBRecursion -= 1
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewFeatured(mod Mod) *Featured {
    f := &Featured{
        Mod: mod,
        ModID: mod.ID,
    }
    f.Meta = "{}"
    return f
}
