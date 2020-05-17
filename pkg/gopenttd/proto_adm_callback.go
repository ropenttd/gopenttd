package gopenttd

import (
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
	"github.com/sirupsen/logrus"
	"sync"
)

// Packet handlers (the mux)

type CallbackMux struct {
	mutex            sync.RWMutex
	onFull           func(packet *packets.ServerFull, conn *OpenttdAdminConnection)
	onBanned         func(packet *packets.ServerBanned, conn *OpenttdAdminConnection)
	onError          func(packet *packets.ServerError, conn *OpenttdAdminConnection)
	onProtocol       func(packet *packets.ServerProtocol, conn *OpenttdAdminConnection)
	onWelcome        func(packet *packets.ServerWelcome, conn *OpenttdAdminConnection)
	onNewgame        func(packet *packets.ServerNewgame, conn *OpenttdAdminConnection)
	onShutdown       func(packet *packets.ServerShutdown, conn *OpenttdAdminConnection)
	onDate           func(packet *packets.ServerDate, conn *OpenttdAdminConnection)
	onClientJoin     func(packet *packets.ServerClientJoin, conn *OpenttdAdminConnection)
	onClientInfo     func(packet *packets.ServerClientInfo, conn *OpenttdAdminConnection)
	onClientUpdate   func(packet *packets.ServerClientUpdate, conn *OpenttdAdminConnection)
	onClientQuit     func(packet *packets.ServerClientQuit, conn *OpenttdAdminConnection)
	onClientError    func(packet *packets.ServerClientError, conn *OpenttdAdminConnection)
	onCompanyNew     func(packet *packets.ServerCompanyNew, conn *OpenttdAdminConnection)
	onCompanyInfo    func(packet *packets.ServerCompanyInfo, conn *OpenttdAdminConnection)
	onCompanyUpdate  func(packet *packets.ServerCompanyUpdate, conn *OpenttdAdminConnection)
	onCompanyRemove  func(packet *packets.ServerCompanyRemove, conn *OpenttdAdminConnection)
	onCompanyEconomy func(packet *packets.ServerCompanyEconomy, conn *OpenttdAdminConnection)
	onCompanyStats   func(packet *packets.ServerCompanyStats, conn *OpenttdAdminConnection)
	onChat           func(packet *packets.ServerChat, conn *OpenttdAdminConnection)
	onRcon           func(packet *packets.ServerRcon, conn *OpenttdAdminConnection)
	onConsole        func(packet *packets.ServerConsole, conn *OpenttdAdminConnection)
	onCmdNames       func(packet *packets.ServerCmdNames, conn *OpenttdAdminConnection)
	onCmdLogging     func(packet *packets.ServerCmdLogging, conn *OpenttdAdminConnection)
	onGamescript     func(packet *packets.ServerGamescript, conn *OpenttdAdminConnection)
	onRconEnd        func(packet *packets.ServerRconEnd, conn *OpenttdAdminConnection)
	onPong           func(packet *packets.ServerPong, conn *OpenttdAdminConnection)
}

func (c *OpenttdAdminConnection) RegFull(f func(packet *packets.ServerFull, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Full handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onFull != nil {
		panic("gopenttd: multiple registrations for ServerFull packet type")
	}

	c.callbacks.onFull = f
}

func (c *OpenttdAdminConnection) RegBanned(f func(packet *packets.ServerBanned, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Banned handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onBanned != nil {
		panic("gopenttd: multiple registrations for ServerBanned packet type")
	}

	c.callbacks.onBanned = f
}

func (c *OpenttdAdminConnection) RegError(f func(packet *packets.ServerError, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Error handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onError != nil {
		panic("gopenttd: multiple registrations for ServerError packet type")
	}

	c.callbacks.onError = f
}

func (c *OpenttdAdminConnection) RegProtocol(f func(packet *packets.ServerProtocol, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Protocol handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onProtocol != nil {
		panic("gopenttd: multiple registrations for ServerProtocol packet type")
	}

	c.callbacks.onProtocol = f
}

func (c *OpenttdAdminConnection) RegWelcome(f func(packet *packets.ServerWelcome, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Welcome handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onWelcome != nil {
		panic("gopenttd: multiple registrations for ServerWelcome packet type")
	}

	c.callbacks.onWelcome = f
}

func (c *OpenttdAdminConnection) RegNewgame(f func(packet *packets.ServerNewgame, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Newgame handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onNewgame != nil {
		panic("gopenttd: multiple registrations for ServerNewgame packet type")
	}

	c.callbacks.onNewgame = f
}

func (c *OpenttdAdminConnection) RegShutdown(f func(packet *packets.ServerShutdown, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Shutdown handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onShutdown != nil {
		panic("gopenttd: multiple registrations for ServerShutdown packet type")
	}

	c.callbacks.onShutdown = f
}

func (c *OpenttdAdminConnection) RegDate(f func(packet *packets.ServerDate, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Date handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onDate != nil {
		panic("gopenttd: multiple registrations for ServerDate packet type")
	}

	c.callbacks.onDate = f
}

func (c *OpenttdAdminConnection) RegClientJoin(f func(packet *packets.ServerClientJoin, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering ClientJoin handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onClientJoin != nil {
		panic("gopenttd: multiple registrations for ServerClientJoin packet type")
	}

	c.callbacks.onClientJoin = f
}

func (c *OpenttdAdminConnection) RegClientInfo(f func(packet *packets.ServerClientInfo, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering ClientInfo handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onClientInfo != nil {
		panic("gopenttd: multiple registrations for ServerClientInfo packet type")
	}

	c.callbacks.onClientInfo = f
}

func (c *OpenttdAdminConnection) RegClientUpdate(f func(packet *packets.ServerClientUpdate, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering ClientUpdate handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onClientUpdate != nil {
		panic("gopenttd: multiple registrations for ServerClientUpdate packet type")
	}

	c.callbacks.onClientUpdate = f
}

func (c *OpenttdAdminConnection) RegClientQuit(f func(packet *packets.ServerClientQuit, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering ClientQuit handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onClientQuit != nil {
		panic("gopenttd: multiple registrations for ServerClientQuit packet type")
	}

	c.callbacks.onClientQuit = f
}

func (c *OpenttdAdminConnection) RegClientError(f func(packet *packets.ServerClientError, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering ClientError handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onClientError != nil {
		panic("gopenttd: multiple registrations for ServerClientError packet type")
	}

	c.callbacks.onClientError = f
}

func (c *OpenttdAdminConnection) RegCompanyNew(f func(packet *packets.ServerCompanyNew, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyNew handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyNew != nil {
		panic("gopenttd: multiple registrations for ServerCompanyNew packet type")
	}

	c.callbacks.onCompanyNew = f
}

func (c *OpenttdAdminConnection) RegCompanyInfo(f func(packet *packets.ServerCompanyInfo, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyInfo handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyInfo != nil {
		panic("gopenttd: multiple registrations for ServerCompanyInfo packet type")
	}

	c.callbacks.onCompanyInfo = f
}

func (c *OpenttdAdminConnection) RegCompanyUpdate(f func(packet *packets.ServerCompanyUpdate, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyUpdate handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyUpdate != nil {
		panic("gopenttd: multiple registrations for ServerCompanyUpdate packet type")
	}

	c.callbacks.onCompanyUpdate = f
}

func (c *OpenttdAdminConnection) RegCompanyRemove(f func(packet *packets.ServerCompanyRemove, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyRemove handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyRemove != nil {
		panic("gopenttd: multiple registrations for ServerCompanyRemove packet type")
	}

	c.callbacks.onCompanyRemove = f
}

func (c *OpenttdAdminConnection) RegCompanyEconomy(f func(packet *packets.ServerCompanyEconomy, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyEconomy handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyEconomy != nil {
		panic("gopenttd: multiple registrations for ServerCompanyEconomy packet type")
	}

	c.callbacks.onCompanyEconomy = f
}

func (c *OpenttdAdminConnection) RegCompanyStats(f func(packet *packets.ServerCompanyStats, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CompanyStats handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCompanyStats != nil {
		panic("gopenttd: multiple registrations for ServerCompanyStats packet type")
	}

	c.callbacks.onCompanyStats = f
}

func (c *OpenttdAdminConnection) RegChat(f func(packet *packets.ServerChat, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Chat handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onChat != nil {
		panic("gopenttd: multiple registrations for ServerChat packet type")
	}

	c.callbacks.onChat = f
}

func (c *OpenttdAdminConnection) RegRcon(f func(packet *packets.ServerRcon, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Rcon handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onRcon != nil {
		panic("gopenttd: multiple registrations for ServerRcon packet type")
	}

	c.callbacks.onRcon = f
}

func (c *OpenttdAdminConnection) RegConsole(f func(packet *packets.ServerConsole, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Console handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onConsole != nil {
		panic("gopenttd: multiple registrations for ServerConsole packet type")
	}

	c.callbacks.onConsole = f
}

func (c *OpenttdAdminConnection) RegCmdNames(f func(packet *packets.ServerCmdNames, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CmdNames handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCmdNames != nil {
		panic("gopenttd: multiple registrations for ServerCmdNames packet type")
	}

	c.callbacks.onCmdNames = f
}

func (c *OpenttdAdminConnection) RegCmdLogging(f func(packet *packets.ServerCmdLogging, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering CmdLogging handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onCmdLogging != nil {
		panic("gopenttd: multiple registrations for ServerCmdLogging packet type")
	}

	c.callbacks.onCmdLogging = f
}

func (c *OpenttdAdminConnection) RegGamescript(f func(packet *packets.ServerGamescript, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Gamescript handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onGamescript != nil {
		panic("gopenttd: multiple registrations for ServerGamescript packet type")
	}

	c.callbacks.onGamescript = f
}

func (c *OpenttdAdminConnection) RegRconEnd(f func(packet *packets.ServerRconEnd, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering RconEnd handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onRconEnd != nil {
		panic("gopenttd: multiple registrations for ServerRconEnd packet type")
	}

	c.callbacks.onRconEnd = f
}

func (c *OpenttdAdminConnection) RegPong(f func(packet *packets.ServerPong, conn *OpenttdAdminConnection)) {
	logrus.Debug("Registering Pong handler")

	c.callbacks.mutex.Lock()
	defer c.callbacks.mutex.Unlock()

	if f == nil {
		panic("gopenttd: nil handler")
	}
	if c.callbacks.onPong != nil {
		panic("gopenttd: multiple registrations for ServerPong packet type")
	}

	c.callbacks.onPong = f
}
