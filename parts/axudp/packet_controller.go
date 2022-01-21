package axudp

//
//type PacketController struct {
//	mandatoryLastId       uint64
//	mandatoryConsistently map[uint64]*packetCollector
//	mandatory             map[uint64]*packetCollector
//
//	optionalLastId       uint64
//	optionalConsistently map[uint64]*packetCollector
//	optional             map[uint64]*packetCollector
//	lock                 sync.RWMutex
//	C                    chan []byte
//}
//
//func NewPacketController() *PacketController {
//	return &PacketController{
//		mandatoryLastId:       0,
//		mandatoryConsistently: map[uint64]*packetCollector{},
//		mandatory:             map[uint64]*packetCollector{},
//		optionalLastId:        0,
//		optionalConsistently:  map[uint64]*packetCollector{},
//		optional:              map[uint64]*packetCollector{},
//		lock:                  sync.RWMutex{},
//		C:                     make(chan []byte),
//	}
//}
//
//func (c *PacketController) ApplyPacket(pck *pproto.Packet) (err error) {
//	var pc *packetCollector
//	switch pck.Type {
//	case pproto.PacketType_PT_OPTIONAL:
//		pc, err = c.addOptionalPacket(pck)
//	default:
//		return fmt.Errorf("packet type %s not implemented", pck.Type.String())
//	}
//	if err != nil {
//		return
//	}
//	if pc.done {
//		log.Info().Uint64("id", pc.id).Str("type", pck.Type.String()).Msg("pck is done")
//		select {
//		case c.C <- pc.bytes():
//		default:
//		}
//	}
//	return
//}
//
//func (c *PacketController) addOptionalPacket(pck *pproto.Packet) (*packetCollector, error) {
//	c.lock.Lock()
//	defer c.lock.Unlock()
//	pc, ok := c.optional[pck.Id]
//	if !ok {
//		pc = newPacketCollector(pck)
//		c.optional[pck.Id] = pc
//	}
//	err := pc.add(pck)
//	if err != nil {
//		return nil, err
//	}
//	return pc, nil
//}
