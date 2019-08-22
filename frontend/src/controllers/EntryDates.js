import moment from "moment";

// Check if a given YYYY-MM-DD date string is valid for an entry.
function isValidEntryDate(d) {
  const m = moment(d);
  if (!m.isValid()) {
    return false;
  }
  const whatGotDoneCreationYear = 2019;
  if (m.year() < whatGotDoneCreationYear) {
    return false;
  }
  if (m > moment(thisFriday())) {
    return false;
  }
  const friday = 5;
  if (m.isoWeekday() != friday) {
    return false;
  }
  return true;
}

// Return the next Friday from the current date in YYYY-MM-DD format.
// If today is Friday, return today's date.
function thisFriday() {
  const today = moment().isoWeekday();
  const friday = 5;

  if (today <= friday) {
    return moment().isoWeekday(friday).format("YYYY-MM-DD");
  } else {
    return moment()
      .add(1, "weeks")
      .isoWeekday(friday).format("YYYY-MM-DD");
  }
}

export { isValidEntryDate, thisFriday };