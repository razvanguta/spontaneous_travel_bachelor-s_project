fetch("/jsonWeather/"+ document.getElementById("theUsername").innerHTML)
.then(response => response.json())
.then(data => {
    var divW = document.getElementById("weather");

    var img = document.createElement("img");
    img.setAttribute("id","icon");    
    var icon = data.weather[0].icon;
    img.src = "http://openweathermap.org/img/w/"+icon+".png"
    divW.appendChild(img)

    const paragraph = document.createElement("p");
    paragraph.setAttribute("id","temp");
    divW.appendChild(paragraph);
    document.getElementById("temp").innerHTML = "Temperatura este de "+data.main.temp+" grade Celsius, iar tempretaura resimtita este " + data.main.feels_like +" grade Celsius";

    const paragraph2 = document.createElement("p");
    paragraph2.setAttribute("id","hum");
    divW.appendChild(paragraph2);
    document.getElementById("hum").innerHTML = "Umiditatea este de " + data.main.humidity;

    const paragraph3 = document.createElement("p");
    paragraph3.setAttribute("id","wind");
    divW.appendChild(paragraph3);
    document.getElementById("wind").innerHTML = "Viteza vantului este de " + data.wind.speed + " km/h";




});
