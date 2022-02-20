fetch("/jsonAgencyBoughtTrips/")
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        document.getElementById("body").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph1 = document.createElement("p");
        paragraph1.setAttribute("id","name-paragraph"+i);
        child.appendChild(paragraph1);
        document.getElementsByTagName("p")[nr].innerHTML=obj.name;
        nr = nr + 1;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","address-paragraph"+i);
        child.appendChild(paragraph2);
        document.getElementsByTagName("p")[nr].innerHTML=obj.address;
        nr = nr + 1;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","email-paragraph"+i);
        child.appendChild(paragraph3);
        document.getElementsByTagName("p")[nr].innerHTML=obj.email;
        nr = nr + 1
        //add the text
        const paragraph4 = document.createElement("p");
        paragraph4.setAttribute("id","phone-paragraph"+i);
        child.appendChild(paragraph4);
        document.getElementsByTagName("p")[nr].innerHTML=obj.phone;
        nr = nr + 1;
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","hotel-paragraph"+i);
        child.appendChild(paragraph5);
        document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel;
        nr = nr + 1;
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","date-paragraph"+i);
        child.appendChild(paragraph6);
        document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date;
        nr = nr + 1;
        //add the text
        const paragraph7 = document.createElement("p");
        paragraph7.setAttribute("id","location-paragraph"+i);
        child.appendChild(paragraph7);
        document.getElementsByTagName("p")[nr].innerHTML="Orasul:-" + obj.city;
        nr = nr + 1;
        //add button
        var button3 = document.createElement("a");
        button3.innerHTML = "Descarca detaliile";
        button3.setAttribute("id","download");
        button3.href = "assets\\pdf\\" + obj.name + obj.city + obj.date + ".pdf";  
        button3.download = "";
        child.appendChild(button3); 
        var separator = document.createElement("hr");
        body.appendChild(separator);
    }
}) 