fetch("/jsonAllTrips/")
.then(response => response.json())
.then(data => {
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        if(obj.agencyName == document.getElementById("theUsername").innerHTML){
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph"+i);
        var body = document.getElementById("body");
        console.log(body)
        body.appendChild(paragraph);
        document.getElementById("agency-paragraph"+i).innerHTML="Agentia de turism "+obj.agencyName;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","title-paragraph"+i);
        body.appendChild(paragraph2);
        document.getElementById("title-paragraph"+i).innerHTML=obj.title;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","desc-paragraph"+i);
        body.appendChild(paragraph3);
        document.getElementById("desc-paragraph"+i).innerHTML=obj.description;
        //add the text
        const paragraph7 = document.createElement("p");
        paragraph7.setAttribute("id","location-paragraph"+i);
        body.appendChild(paragraph7);
        document.getElementById("location-paragraph"+i).innerHTML="Orasul: " + obj.city + " Tara: " + obj.country;
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","hotel-paragraph"+i);
        body.appendChild(paragraph6);
        document.getElementById("hotel-paragraph"+i).innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
        //add the text
        const paragraph4 = document.createElement("p");
        paragraph4.setAttribute("id","price-paragraph"+i);
        body.appendChild(paragraph4);
        document.getElementById("price-paragraph"+i).innerHTML=obj.price+" euro";
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","details-paragraph"+i);
        body.appendChild(paragraph5);
        document.getElementById("details-paragraph"+i).innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
        //add the photo
        var img1 = document.createElement("img");
        img1.setAttribute("id","agency-photo");        
        img1.src = "../"+obj.img1
        body.appendChild(img1)
        //add the photo
        var img2 = document.createElement("img");
        img2.setAttribute("id","agency-photo");        
        img2.src = "../"+obj.img2
        body.appendChild(img2)
        //add the photo
        var img3 = document.createElement("img");
        img3.setAttribute("id","agency-photo");        
        img3.src = "../"+obj.img3
        body.appendChild(img3)
        //add space
        var space = document.createElement("br");
        var separator = document.createElement("hr");
        body.appendChild(space);


        //add a button
        if(obj.same == "yes"){
            var form1 = document.createElement("form");
            form1.setAttribute("id","deleteTrip");
            form1.method="GET";
            form1.action="/deleteTrip/"+obj.id;
            body.appendChild(form1);
            var button = document.createElement("button");
            button.innerHTML = "Sterge excursia";
            button.setAttribute("id","delte-button");
            button.type="submit";
            form1.appendChild(button); 
            body.append(space);

            var form2 = document.createElement("form");
            form2.setAttribute("id","updateTrip");
            form2.method="GET";
            form2.action="/updateTripPage/"+obj.id;
            body.appendChild(form2);
            var button2 = document.createElement("button");
            button2.innerHTML = "Editeaza excursia";
            button2.setAttribute("id","update-button");
            button2.type="submit";
            form2.appendChild(button2); 
            body.append(space)
            }

            var form3 = document.createElement("form");
            form3.setAttribute("id","getWeather");
            form3.method="GET";
            form3.action="/weather/"+obj.city;  
            body.appendChild(form3);
            var button3 = document.createElement("button");
            button3.innerHTML = "Vremea in "+ obj.city;
            button3.type="submit";
            button3.setAttribute("id","weather");
            form3.appendChild(button3); 


        }


    }
}) 