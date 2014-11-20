// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webgo

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// SessionManager manager sessions.
type SessionManager struct {
	mutex      sync.RWMutex
	mutexOnEnd sync.RWMutex
	sessionMap map[string]*Session
	onStart    func(*Session)
	onTouch    func(*Session)
	onEnd      func(*Session)
	path       string
	timeout    uint
}

// Session stores the id and value for a session.
type Session struct {
	Id    string
	Value interface{}

	manager *SessionManager
	res     http.ResponseWriter
	expire  int64
}

func (session *Session) Cookie() string {
	tm := time.Unix(session.expire, 0).UTC()
	return fmt.Sprintf(
		"SessionId=%s; path=%s; expires=%s;",
		session.Id,
		session.manager.path,
		tm.Format("Fri, 02-Jan-2006 15:04:05 -0700"),
	)
}

func (session *Session) Abandon() {
	if _, found := session.manager.sessionMap[session.Id]; found {
		delete((*session.manager).sessionMap, session.Id)
	}
	if session.res != nil {
		session.res.Header().Set(
			"Set-Cookie", fmt.Sprintf("SessionId=; path=%s;", session.manager.path),
		)
	}
}

func NewSessionManager(logger *log.Logger) *SessionManager {
	p := &SessionManager{
		path:       "/",
		sessionMap: make(map[string]*Session),
		timeout:    300,
	}
	go func(p *SessionManager) {
		for { // never stop !!!
			l := time.Now().Unix()
			if f := p.onEnd; f != nil {
				p.mutexOnEnd.Lock()
				for id, v := range p.sessionMap {
					if v.expire < l {
						if logger != nil {
							logger.Printf("Expired session(id:%s)", id)
						}
						f(v)
						delete(p.sessionMap, id)
					}
				}
				p.mutexOnEnd.Unlock()
			} else {
				for id, v := range p.sessionMap {
					if v.expire < l {
						if logger != nil {
							logger.Printf("Expired session(id:%s)", id)
						}
						delete(p.sessionMap, id)
					}
				}
			}
			time.Sleep(time.Second)
		}
	}(p)
	return p
}

func (p *SessionManager) Has(id string) (found bool) {
	_, found = p.sessionMap[id]
	return
}

func (p *SessionManager) GetSessionById(id string) (session *Session) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if id == "" || !p.Has(id) {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return
		}
		id = fmt.Sprintf("%x", b)
	}
	tm := time.Unix(time.Now().Unix()+int64(p.timeout), 0).UTC()
	var found bool
	if session, found = p.sessionMap[id]; found {
		session.expire = tm.Unix()
		if f := p.onTouch; f != nil {
			f(session)
		}
		return
	} else {
		session = &Session{Id: id, expire: tm.Unix(), manager: p}
		p.sessionMap[id] = session
		if f := p.onStart; f != nil {
			f(session)
		}
	}
	return
}

func (p *SessionManager) GetSession(res http.ResponseWriter, req *http.Request) (session *Session) {
	if c, _ := req.Cookie("SessionId"); c != nil {
		session = p.GetSessionById(c.Value)
	} else {
		session = p.GetSessionById("")
	}
	if res != nil {
		session.res = res
		res.Header().Add("Set-Cookie",
			fmt.Sprintf("SessionId=%s; path=%s; expires=%s;",
				session.Id,
				session.manager.path,
				time.Unix(session.expire, 0).UTC().Format(
					"Fri, 02-Jan-2006 15:04:05 GMT",
				),
			),
		)
	}
	return
}

func (p *SessionManager) Abandon() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if f := p.onEnd; f != nil {
		p.mutexOnEnd.Lock()
		for id, v := range p.sessionMap {
			f(v)
			delete(p.sessionMap, id)
		}
		p.mutexOnEnd.Unlock()
	} else {
		for id, _ := range p.sessionMap {
			delete(p.sessionMap, id)
		}
	}
}

func (p *SessionManager) OnStart(f func(*Session)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.onStart = f
}

func (p *SessionManager) OnTouch(f func(*Session)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.onTouch = f
}

func (p *SessionManager) OnEnd(f func(*Session)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.onEnd = f
}

func (p *SessionManager) SetTimeout(t uint) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.timeout = t
}

func (p *SessionManager) GetTimeout() uint {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.timeout
}

func (p *SessionManager) SetPath(t string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.path = t
}

func (p *SessionManager) GetPath() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.path
}
