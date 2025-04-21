package ziface

type IDataPack interface {
	// 获取包的头长度
	GetHeadLen() uint32
	// 封包方法：将消息数据封装成一个Message
	Pack(IMessage) ([]byte, error)
	// 拆包方法：将数据拆分到Message中
	UnPack([]byte) (IMessage, error)
}
