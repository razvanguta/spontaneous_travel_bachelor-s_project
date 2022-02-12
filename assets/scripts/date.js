var date = new Date();
//the trip can be added not sooner than 2 days
date.setDate(date.getDate() + 2);
var startDate = (date).toISOString().substring(0,10);
document.getElementById("date").min = startDate;
//the trip can be added not later than 52 days(it's an last minute)
date.setDate(date.getDate() + 50);
var endDate = (date).toISOString().substring(0,10);
document.getElementById("date").max = endDate;