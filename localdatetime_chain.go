package goda

import "time"

type LocalDateTimeChain struct {
	Chain[LocalDateTime]
}

func (l LocalDateTimeChain) PlusYears(years int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusYears"))
	l.value.date = l.value.date.chainWithError(l.eError).PlusYears(years).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) MinusYears(years int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusYears"))
	l.value.date = l.value.date.chainWithError(l.eError).MinusYears(years).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) PlusMonths(months int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusMonths"))
	l.value.date = l.value.date.chainWithError(l.eError).PlusMonths(months).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) MinusMonths(months int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusMonths"))
	l.value.date = l.value.date.chainWithError(l.eError).MinusMonths(months).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) PlusWeeks(weeks int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusWeeks"))
	l.value.date = l.value.date.chainWithError(l.eError).PlusWeeks(weeks).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) MinusWeeks(weeks int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusWeeks"))
	l.value.date = l.value.date.chainWithError(l.eError).MinusWeeks(weeks).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) PlusDays(days int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusDays"))
	l.value.date = l.value.date.chainWithError(l.eError).PlusDays(days).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) MinusDays(days int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusDays"))
	l.value.date = l.value.date.chainWithError(l.eError).MinusDays(days).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) PlusHours(hours int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusHours"))
	return l.plusWithOverflow(l.value.date, hours, 0, 0, 0, 1)
}

func (l LocalDateTimeChain) MinusHours(hours int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusHours"))
	return l.plusWithOverflow(l.value.date, hours, 0, 0, 0, -1)
}

func (l LocalDateTimeChain) PlusMinutes(minutes int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusMinutes"))
	return l.plusWithOverflow(l.value.date, 0, minutes, 0, 0, 1)
}

func (l LocalDateTimeChain) MinusMinutes(minutes int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusMinutes"))
	return l.plusWithOverflow(l.value.date, 0, minutes, 0, 0, -1)
}

func (l LocalDateTimeChain) PlusSeconds(seconds int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusSeconds"))
	return l.plusWithOverflow(l.value.date, 0, 0, seconds, 0, 1)
}

func (l LocalDateTimeChain) MinusSeconds(seconds int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusSeconds"))
	return l.plusWithOverflow(l.value.date, 0, 0, seconds, 0, -1)
}

func (l LocalDateTimeChain) PlusNanos(nanos int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "PlusNanos"))
	return l.plusWithOverflow(l.value.date, 0, 0, 0, nanos, 1)
}

func (l LocalDateTimeChain) MinusNanos(nanos int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "MinusNanos"))
	return l.plusWithOverflow(l.value.date, 0, 0, 0, nanos, -1)
}

func (l LocalDateTimeChain) plusWithOverflow(newDate LocalDate, hours, minutes, seconds, nanos, sign int64) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "plusWithOverflow"))
	if !l.ok() {
		return l
	}
	if hours|minutes|seconds|nanos == 0 {
		l.value.date = newDate
		return l
	}
	const (
		NanosPerSecond = int64(time.Second)
		NanosPerMinute = int64(time.Minute)
		NanosPerHour   = int64(time.Hour)
		NanosPerDay    = int64(time.Hour * 24)

		SecondsPerDay = int64(24 * 60 * 60)
		MinutesPerDay = int64(24 * 60)
		HoursPerDay   = int64(24)
	)
	var totDays = nanos/NanosPerDay +
		seconds/SecondsPerDay +
		minutes/MinutesPerDay +
		hours/HoursPerDay
	totDays *= sign
	var totNanos = nanos%NanosPerDay +
		(seconds%SecondsPerDay)*NanosPerSecond +
		(minutes%MinutesPerDay)*NanosPerMinute +
		(hours%HoursPerDay)*NanosPerHour
	var curNoD = l.value.time.NanoOfDay()
	totNanos = totNanos*sign + curNoD
	totDays += floorDiv(totNanos, NanosPerDay)
	var newNoD = floorMod(totNanos, NanosPerDay)
	if newNoD != curNoD {
		l.value.time, l.eError = LocalTimeOfNanoOfDay(newNoD)
	}
	l.value.date = newDate.Chain().PlusDays(totDays).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithYear(year Year) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithYear"))
	l.value.date = l.value.date.chainWithError(l.eError).WithYear(year).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithMonth(month Month) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithMonth"))
	l.value.date = l.value.date.chainWithError(l.eError).WithMonth(month).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithDayOfMonth(dayOfMonth int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithDayOfMonth"))
	l.value.date = l.value.date.chainWithError(l.eError).WithDayOfMonth(dayOfMonth).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithDayOfYear(dayOfYear int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithDayOfYear"))
	l.value.date = l.value.date.chainWithError(l.eError).WithDayOfYear(dayOfYear).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithHour(hour int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithHour"))
	l.value.time = l.value.time.chainWithError(l.eError).WithHour(hour).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithMinute(minute int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithMinute"))
	l.value.time = l.value.time.chainWithError(l.eError).WithMinute(minute).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithSecond(second int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithSecond"))
	l.value.time = l.value.time.chainWithError(l.eError).WithSecond(second).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithNano(nanoOfSecond int) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithNano"))
	l.value.time = l.value.time.chainWithError(l.eError).WithNano(nanoOfSecond).mergeError(&l.eError)
	return l
}

func (l LocalDateTimeChain) WithField(field Field, value TemporalValue) LocalDateTimeChain {
	defer l.leaveFunction(l.enterFunction("LocalDateTime", "WithField"))
	if field.IsTimeBased() {
		l.value.time = l.value.time.chainWithError(l.eError).WithField(field, value).mergeError(&l.eError)
	} else {
		l.value.date = l.value.date.chainWithError(l.eError).WithField(field, value).mergeError(&l.eError)
	}
	return l
}
