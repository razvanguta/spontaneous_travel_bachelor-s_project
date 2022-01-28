var date = new Date();
var currentDate = date.toISOString().substring(0,10);
console.log(currentDate);
document.getElementById("date").min = currentDate;