fetch("/jsonClientBoughtTrips/")
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        const divAll = document.createElement("div");
        divAll.setAttribute("id","div"+i);
        divAll.setAttribute("class","formular-login2")
        document.getElementById("body").appendChild(divAll);
        var child = document.getElementById("div"+i);
        //add the text
        const paragraph2 = document.createElement("h2");
        paragraph2.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph2);
        paragraph2.innerHTML=obj.title;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","desc-paragraph"+i);
        child.appendChild(paragraph3);
        document.getElementsByTagName("p")[nr].innerHTML=obj.description;
        nr = nr + 1;
        //add the text
        const paragraph7 = document.createElement("p");
        paragraph7.setAttribute("id","location-paragraph"+i);
        child.appendChild(paragraph7);
        document.getElementsByTagName("p")[nr].innerHTML="Orasul:-" + obj.city + "-Tara: " + obj.country;
        nr = nr + 1
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","hotel-paragraph"+i);
        child.appendChild(paragraph6);
        document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel;
        nr = nr + 1;
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","date-paragraph"+i);
        child.appendChild(paragraph5);
        document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date;
        nr = nr + 1;
         //add button
         var form = document.createElement("form");
         form.setAttribute("id","downloadForm");
         form.method="get";
         form.action="assets\\pdf\\" + obj.clientId + obj.city + obj.date + ".pdf";
         var button3 = document.createElement("button");
         button3.innerHTML = "Descarca detaliile";
         button3.setAttribute("id","download");
         button3.setAttribute("class","toRegister-btn");
         button3.type = "submit";
         form.appendChild(button3); 
         child.appendChild(form);
        //add space
        var separator = document.createElement("hr");
        body.appendChild(separator);
    }
}) 