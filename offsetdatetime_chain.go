package goda

type OffsetDateTimeChain struct {
	Chain[OffsetDateTime]
}

func (o OffsetDateTimeChain) WithField(field Field, value TemporalValue) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithField"))
	newValue := value.Int64()
	field.checkSetE(newValue, &o.eError)
	if !o.ok() {
		return o
	}
	switch field {
	case FieldInstantSeconds:
		o.value.datetime, o.eError = LocalDateTimeOfEpochSecond(newValue+int64(o.value.offset.totalSeconds), int64(o.value.Nanosecond()), o.value.offset)
	case FieldOffsetSeconds:
		o.value.offset.totalSeconds = int32(newValue)
	default:
		o.value.datetime = o.value.datetime.chainWithError(o.eError).WithField(field, value).mergeError(&o.eError)
	}
	return o
}

func (o OffsetDateTimeChain) PlusYears(years int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusYears"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusYears(years).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusYears(years int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusYears"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusYears(years).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusMonths(months int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusMonths"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusMonths(months).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusMonths(months int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusMonths"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusMonths(months).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusWeeks(weeks int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusWeeks"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusWeeks(weeks).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusWeeks(weeks int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusWeeks"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusWeeks(weeks).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusDays(days int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusDays"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusDays(days).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusDays(days int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusDays"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusDays(days).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusHours(hours int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusHours"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusHours(hours).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusHours(hours int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusHours"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusHours(hours).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusMinutes(minutes int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusMinutes"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusMinutes(minutes).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusMinutes(minutes int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusMinutes"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusMinutes(minutes).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusSeconds(seconds int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusSeconds"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusSeconds(seconds).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusSeconds(seconds int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusSeconds"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusSeconds(seconds).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) PlusNanos(nanos int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "PlusNanos"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).PlusNanos(nanos).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) MinusNanos(nanos int64) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "MinusNanos"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).MinusNanos(nanos).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithYear(year Year) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithYear"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithYear(year).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithMonth(month Month) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithMonth"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithMonth(month).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithDayOfMonth(dayOfMonth int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithDayOfMonth"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithDayOfMonth(dayOfMonth).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithDayOfYear(dayOfYear int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithDayOfYear"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithDayOfYear(dayOfYear).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithHour(hour int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithHour"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithHour(hour).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithMinute(minute int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithMinute"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithMinute(minute).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithSecond(second int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithSecond"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithSecond(second).mergeError(&o.eError)
	return o
}
func (o OffsetDateTimeChain) WithNano(nanoOfSecond int) OffsetDateTimeChain {
	defer o.leaveFunction(o.enterFunction("OffsetDateTime", "WithNano"))
	o.value.datetime = o.value.datetime.chainWithError(o.eError).WithNano(nanoOfSecond).mergeError(&o.eError)
	return o
}
