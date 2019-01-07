package accessories

import "errors"

type CustomSerialController struct {
	comm *SerialAccessory
}

func (csc *CustomSerialController) CheckConnectionOK() error {
	resp, err := csc.comm.SendCommand([]byte("ok?\r"))
	if err != nil {
		return err
	}
	if string(resp) != ">ok?:ok\r" {
		return errors.New("serial device didn't respond correctly")
	}
	return nil
}

func (csc *CustomSerialController) UnlockAllDoors() error {
	resp, err := csc.comm.SendCommand([]byte("body/unlock/all\r"))
	if err != nil {
		return err
	}
	if string(resp) != ">body/unlock/all:ok\r" {
		return errors.New("serial device didn't respond correctly")
	}
	return nil
}

func (csc *CustomSerialController) LockAllDoors() error {
	resp, err := csc.comm.SendCommand([]byte("body/lock/all\r"))
	if err != nil {
		return err
	}
	if string(resp) != ">body/lock/all:ok\r" {
		return errors.New("serial device didn't respond correctly")
	}
	return nil
}
