package goda

import "time"

type LocalDateTimeChain struct {
	Chain[LocalDateTime]
}

func (l LocalDateTimeChain) PlusYears(years int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).PlusYears(years).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) MinusYears(years int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).MinusYears(years).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) PlusMonths(months int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).PlusMonths(months).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) MinusMonths(months int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).MinusMonths(months).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) PlusWeeks(weeks int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).PlusWeeks(weeks).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) MinusWeeks(weeks int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).MinusWeeks(weeks).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) PlusDays(days int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).PlusDays(days).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) MinusDays(days int64) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).MinusDays(days).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) PlusHours(hours int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, hours, 0, 0, 0, 1)
}

func (l LocalDateTimeChain) MinusHours(hours int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, -hours, 0, 0, 0, -1)
}

func (l LocalDateTimeChain) PlusMinutes(minutes int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, minutes, 0, 0, 1)
}

func (l LocalDateTimeChain) MinusMinutes(minutes int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, -minutes, 0, 0, -1)
}

func (l LocalDateTimeChain) PlusSeconds(seconds int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, 0, seconds, 0, 1)
}

func (l LocalDateTimeChain) MinusSeconds(seconds int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, 0, -seconds, 0, -1)
}

func (l LocalDateTimeChain) PlusNanos(nanos int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, 0, 0, nanos, 1)
}

func (l LocalDateTimeChain) MinusNanos(nanos int64) LocalDateTimeChain {
	return l.plusWithOverflow(l.value.date, 0, 0, 0, -nanos, -1)
}

func (l LocalDateTimeChain) plusWithOverflow(newDate LocalDate, hours, minutes, seconds, nanos, sign int64) LocalDateTimeChain {
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
		l.value.time, l.Error = LocalTimeOfNanoOfDay(newNoD)
	}
	l.value.date = newDate.Chain().PlusDays(totDays).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithYear(year Year) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).WithYear(year).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithMonth(month Month) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).WithMonth(month).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithDayOfMonth(dayOfMonth int) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).WithDayOfMonth(dayOfMonth).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithDayOfYear(dayOfYear int) LocalDateTimeChain {
	l.value.date = l.value.date.chainWithError(l.Error).WithDayOfYear(dayOfYear).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithHour(hour int) LocalDateTimeChain {
	l.value.time = l.value.time.chainWithError(l.Error).WithHour(hour).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithMinute(minute int) LocalDateTimeChain {
	l.value.time = l.value.time.chainWithError(l.Error).WithMinute(minute).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithSecond(second int) LocalDateTimeChain {
	l.value.time = l.value.time.chainWithError(l.Error).WithSecond(second).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithNano(nanoOfSecond int) LocalDateTimeChain {
	l.value.time = l.value.time.chainWithError(l.Error).WithNano(nanoOfSecond).getError(&l.Error)
	return l
}

func (l LocalDateTimeChain) WithField(field Field, value TemporalValue) LocalDateTimeChain {
	if field.IsTimeBased() {
		l.value.time = l.value.time.chainWithError(l.Error).WithField(field, value).getError(&l.Error)
	} else {
		l.value.date = l.value.date.chainWithError(l.Error).WithField(field, value).getError(&l.Error)
	}
	return l
}
