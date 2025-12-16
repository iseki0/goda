package goda

type LocalTimeChain struct {
	Chain[LocalTime]
}

func (l LocalTimeChain) PlusHours(hours int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "PlusHours"))
	if hours == 0 {
		return l
	}
	if !l.ok() {
		return l
	}
	newHour := (hours%24 + int64(l.value.Hour()) + 24) % 24
	l.value, l.eError = LocalTimeOf(int(newHour), l.value.Minute(), l.value.Second(), l.value.Nano())
	return l
}

func (l LocalTimeChain) MinusHours(hoursToSubtract int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "MinusHours"))
	return l.PlusHours(-(hoursToSubtract % 24))
}

func (l LocalTimeChain) PlusMinutes(minutesToAdd int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "PlusMinutes"))
	if minutesToAdd == 0 {
		return l
	}
	var mofd = int64(l.value.Hour()*60 + l.value.Minute())
	var newMofd = (minutesToAdd%(60*24) + mofd + (60 * 24)) % (60 * 24)
	if mofd == newMofd {
		return l
	}
	l.value, l.eError = LocalTimeOf(int(newMofd/60), int(newMofd%60), l.value.Second(), l.value.Nano())
	return l
}

func (l LocalTimeChain) MinusMinutes(minutesToSubtract int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "MinusMinutes"))
	return l.PlusMinutes(-(minutesToSubtract % 1440))
}

func (l LocalTimeChain) PlusSeconds(secondsToAdd int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "PlusSeconds"))
	if secondsToAdd == 0 {
		return l
	}
	var sofd = int64(l.value.SecondOfDay())
	var newSofd = (secondsToAdd%86400 + sofd + 86400) % 86400
	if sofd == newSofd {
		return l
	}
	l.value, l.eError = LocalTimeOf(int(newSofd/3600), int(newSofd/60%60), int(newSofd%60), l.value.Nano())
	return l
}

func (l LocalTimeChain) MinusSeconds(secondsToSubtract int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "MinusSeconds"))
	return l.PlusSeconds(-(secondsToSubtract % 86400))
}

func (l LocalTimeChain) PlusNano(nanosToAdd int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "PlusNano"))
	if nanosToAdd == 0 {
		return l
	}
	const NanosPerDay = 86400_000_000_000
	var nofd = l.value.NanoOfDay()
	var newNofd = ((nanosToAdd % NanosPerDay) + nofd + NanosPerDay) % NanosPerDay
	if nofd == newNofd {
		return l
	}
	const NanosPerHour = 3600_000_000_000
	const NanosPerMinute = 60_000_000_000
	var newHour = int(newNofd / NanosPerHour)
	var newMinute = int((newNofd / NanosPerMinute) % 60)
	var newSecond = int((newNofd / 1000_000_000) % 60)
	var newNano = int(newNofd % 1000_000_000)
	l.value, l.eError = LocalTimeOf(newHour, newMinute, newSecond, newNano)
	return l
}

func (l LocalTimeChain) MinusNano(nanosToSubtract int64) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "MinusNano"))
	return l.PlusNano(-(nanosToSubtract % 86400_000_000_000))
}

func (l LocalTimeChain) WithNano(nanoOfSecond int) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "WithNano"))
	FieldNanoOfSecond.checkSetE(int64(nanoOfSecond), &l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalTimeOf(l.value.Hour(), l.value.Minute(), l.value.Second(), nanoOfSecond)
	return l
}

func (l LocalTimeChain) WithSecond(second int) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "WithSecond"))
	FieldSecondOfMinute.checkSetE(int64(second), &l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalTimeOf(l.value.Hour(), l.value.Minute(), second, l.value.Nano())
	return l
}

func (l LocalTimeChain) WithMinute(minute int) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "WithMinute"))
	FieldMinuteOfHour.checkSetE(int64(minute), &l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalTimeOf(l.value.Hour(), minute, l.value.Second(), l.value.Nano())
	return l
}

func (l LocalTimeChain) WithHour(hour int) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "WithHour"))
	FieldHourOfDay.checkSetE(int64(hour), &l.eError)
	if !l.ok() {
		return l
	}
	l.value, l.eError = LocalTimeOf(hour, l.value.Minute(), l.value.Second(), l.value.Nano())
	return l
}

// WithField returns a copy of this LocalTime with the specified field replaced.
// Zero values return zero immediately.
//
// Supported fields mirror Java's LocalTime#with(TemporalField, long):
//   - FieldNanoOfSecond: sets the nano-of-second while keeping hour, minute, and second.
//   - FieldNanoOfDay: replaces the entire time using the provided nano-of-day (equivalent to LocalTimeOfNanoOfDay).
//   - FieldMicroOfSecond: replaces the nano-of-second with micro-of-second × 1,000; hour, minute, and second stay the same.
//   - FieldMicroOfDay: replaces the entire time using micro-of-day × 1,000 (equivalent to LocalTimeOfNanoOfDay).
//   - FieldMilliOfSecond: replaces the nano-of-second with milli-of-second × 1,000,000; hour, minute, and second stay the same.
//   - FieldMilliOfDay: replaces the entire time using milli-of-day × 1,000,000 (equivalent to LocalTimeOfNanoOfDay).
//   - FieldSecondOfMinute: sets the second-of-minute while leaving hour, minute, and nano-of-second untouched.
//   - FieldSecondOfDay: replaces hour, minute, and second based on the provided second-of-day while keeping nano-of-second.
//   - FieldMinuteOfHour: sets the minute-of-hour; hour, second, and nano-of-second stay unchanged.
//   - FieldMinuteOfDay: replaces hour and minute based on the minute-of-day; second and nano-of-second stay unchanged.
//   - FieldHourOfAmPm: sets the hour within AM/PM while keeping the current half of day, minute, second, and nano-of-second.
//   - FieldClockHourOfAmPm: sets the 1-12 clock hour within AM/PM while keeping the current half of day, minute, second, and nano-of-second.
//   - FieldHourOfDay: sets the hour-of-day while leaving minute, second, and nano-of-second untouched.
//   - FieldClockHourOfDay: sets the 1-24 clock hour-of-day (24 → 0) while leaving minute, second, and nano-of-second untouched.
//   - FieldAmPmOfDay: toggles AM/PM while preserving the hour-of-am-pm, minute, second, and nano-of-second.
//
// Fields outside this list return an error. Range violations propagate the validation error.
func (l LocalTimeChain) WithField(field Field, value TemporalValue) LocalTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalTime", "WithField"))
	field.checkSetE(value.Int64(), &l.eError)
	if !l.ok() {
		return l
	}
	newValue := value.v
	var hour = l.value.Hour()
	var minute = l.value.Minute()
	switch field {
	case FieldNanoOfDay:
		l.value, l.eError = LocalTimeOfNanoOfDay(newValue)
	case FieldMicroOfDay:
		l.value, l.eError = LocalTimeOfNanoOfDay(newValue * 1000)
	case FieldMilliOfDay:
		l.value, l.eError = LocalTimeOfNanoOfDay(newValue * 1000_000)
	case FieldNanoOfSecond:
		return l.WithNano(int(newValue))
	case FieldMicroOfSecond:
		return l.WithNano(int(newValue) * 1000)
	case FieldMilliOfSecond:
		return l.WithNano(int(newValue) * 1000_000)
	case FieldSecondOfMinute:
		return l.WithSecond(int(newValue))
	case FieldSecondOfDay:
		return l.PlusSeconds(newValue - int64(l.value.SecondOfDay()))
	case FieldMinuteOfHour:
		return l.WithMinute(int(newValue))
	case FieldMinuteOfDay:
		return l.PlusMinutes(newValue - int64(hour*60+minute))
	case FieldHourOfAmPm:
		if newValue == 12 {
			return l.PlusHours(int64(0 - hour%12))
		}
		return l.PlusHours(newValue - int64(hour%12))
	case FieldHourOfDay:
		return l.WithHour(int(newValue))
	case FieldClockHourOfDay:
		if newValue == 24 {
			return l.WithHour(0)
		}
		return l.WithHour(int(newValue))
	case FieldAmPmOfDay:
		return l.PlusHours((newValue - int64(hour)/12) * 12)
	case FieldClockHourOfAmPm:
		var a = newValue
		if newValue == 12 {
			a = 0
		}
		return l.PlusHours(a - int64(l.value.Hour())%12)
	default:
		l.eError = unsupportedField(field)
	}
	return l
}
