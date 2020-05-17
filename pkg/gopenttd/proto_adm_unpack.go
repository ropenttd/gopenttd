package gopenttd

import (
	"bytes"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
)

func (c *OpenttdAdminConnection) handlePacket(data []byte) (packet packets.AdminResponsePacket, err error) {
	c.callbacks.mutex.RLock()
	defer c.callbacks.mutex.RUnlock()

	reader := bytes.NewBuffer(data)
	packetType := uint8(reader.Next(1)[0])
	switch packetType {
	// there's probably a better way to do this, if you can tidy it please do and PR it!
	case uint8(admin.PacketServerFull):
		obj := packets.ServerFull{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onFull != nil {
			c.callbacks.onFull(&obj, c)
		}
	case uint8(admin.PacketServerBanned):
		obj := packets.ServerBanned{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onBanned != nil {
			c.callbacks.onBanned(&obj, c)
		}
	case uint8(admin.PacketServerError):
		obj := packets.ServerError{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onError != nil {
			c.callbacks.onError(&obj, c)
		}
	case uint8(admin.PacketServerProtocol):
		obj := packets.ServerProtocol{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onProtocol != nil {
			c.callbacks.onProtocol(&obj, c)
		}
	case uint8(admin.PacketServerWelcome):
		obj := packets.ServerWelcome{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onWelcome != nil {
			c.callbacks.onWelcome(&obj, c)
		}
	case uint8(admin.PacketServerNewgame):
		obj := packets.ServerNewgame{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onNewgame != nil {
			c.callbacks.onNewgame(&obj, c)
		}
	case uint8(admin.PacketServerShutdown):
		obj := packets.ServerShutdown{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onShutdown != nil {
			c.callbacks.onShutdown(&obj, c)
		}
	case uint8(admin.PacketServerDate):
		obj := packets.ServerDate{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onDate != nil {
			c.callbacks.onDate(&obj, c)
		}
	case uint8(admin.PacketServerClientJoin):
		obj := packets.ServerClientJoin{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onClientJoin != nil {
			c.callbacks.onClientJoin(&obj, c)
		}
	case uint8(admin.PacketServerClientInfo):
		obj := packets.ServerClientInfo{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onClientInfo != nil {
			c.callbacks.onClientInfo(&obj, c)
		}
	case uint8(admin.PacketServerClientUpdate):
		obj := packets.ServerClientUpdate{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onClientUpdate != nil {
			c.callbacks.onClientUpdate(&obj, c)
		}
	case uint8(admin.PacketServerClientQuit):
		obj := packets.ServerClientQuit{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onClientQuit != nil {
			c.callbacks.onClientQuit(&obj, c)
		}
	case uint8(admin.PacketServerClientError):
		obj := packets.ServerClientError{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onClientError != nil {
			c.callbacks.onClientError(&obj, c)
		}
	case uint8(admin.PacketServerCompanyNew):
		obj := packets.ServerCompanyNew{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyNew != nil {
			c.callbacks.onCompanyNew(&obj, c)
		}
	case uint8(admin.PacketServerCompanyInfo):
		obj := packets.ServerCompanyInfo{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyInfo != nil {
			c.callbacks.onCompanyInfo(&obj, c)
		}
	case uint8(admin.PacketServerCompanyUpdate):
		obj := packets.ServerCompanyUpdate{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyUpdate != nil {
			c.callbacks.onCompanyUpdate(&obj, c)
		}
	case uint8(admin.PacketServerCompanyRemove):
		obj := packets.ServerCompanyRemove{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyRemove != nil {
			c.callbacks.onCompanyRemove(&obj, c)
		}
	case uint8(admin.PacketServerCompanyEconomy):
		obj := packets.ServerCompanyEconomy{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyEconomy != nil {
			c.callbacks.onCompanyEconomy(&obj, c)
		}
	case uint8(admin.PacketServerCompanyStats):
		obj := packets.ServerCompanyStats{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCompanyStats != nil {
			c.callbacks.onCompanyStats(&obj, c)
		}
	case uint8(admin.PacketServerChat):
		obj := packets.ServerChat{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onChat != nil {
			c.callbacks.onChat(&obj, c)
		}
	case uint8(admin.PacketServerRcon):
		obj := packets.ServerRcon{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onRcon != nil {
			c.callbacks.onRcon(&obj, c)
		}
	case uint8(admin.PacketServerCmdNames):
		obj := packets.ServerCmdNames{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCmdNames != nil {
			c.callbacks.onCmdNames(&obj, c)
		}
	case uint8(admin.PacketServerCmdLogging):
		obj := packets.ServerCmdLogging{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onCmdLogging != nil {
			c.callbacks.onCmdLogging(&obj, c)
		}
	case uint8(admin.PacketServerGamescript):
		obj := packets.ServerGamescript{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onGamescript != nil {
			c.callbacks.onGamescript(&obj, c)
		}
	case uint8(admin.PacketServerRconEnd):
		obj := packets.ServerRconEnd{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onRconEnd != nil {
			c.callbacks.onRconEnd(&obj, c)
		}
	case uint8(admin.PacketServerConsole):
		obj := packets.ServerConsole{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onConsole != nil {
			c.callbacks.onConsole(&obj, c)
		}
	case uint8(admin.PacketServerPong):
		obj := packets.ServerPong{}
		err = obj.Unmarshal(reader)
		packet = &obj
		if c.callbacks.onPong != nil {
			c.callbacks.onPong(&obj, c)
		}
	}
	return packet, err
}
