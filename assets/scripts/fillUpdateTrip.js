
fetch("/jsonUpdateTrip/"+ document.getElementById("theID").innerHTML)
.then(response => response.json())
.then(data => {
    document.getElementById("title").value = data.title;
    document.getElementById("description").value =  data.description;
    document.getElementById("hotel").value =  data.hotel;
    document.getElementById("stars").value =  data.stars;
    document.getElementById("price").value =  data.price;
    document.getElementById("date").value =  data.date;
    document.getElementById("days").value =  data.days;
    document.getElementById("city").value = data.city;
    document.getElementById("country").value = data.country
    document.getElementById("form_update").action = "/updateTripLogic/" + document.getElementById("theID").innerHTML
});
