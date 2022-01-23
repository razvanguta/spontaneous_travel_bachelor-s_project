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
            var button = document.createElement("a");
            button.innerHTML = "Sterge excursia";
            button.setAttribute("id","delte-button");
            button.href = "/deleteTrip/"+obj.id;  
            body.appendChild(button); 
            body.append(space)

            var button2 = document.createElement("a");
            button2.innerHTML = "Editeaza excursia";
            button2.setAttribute("id","update-button");
            button2.href = "/UpdateTripPage/"+obj.id;  
            body.appendChild(button2); 
            body.append(space)
            }

            var button3 = document.createElement("a");
            button3.innerHTML = "Vremea in "+ obj.city;
            button3.setAttribute("id","weather");
            button3.href = "/weather/"+obj.city;  
            body.appendChild(button3);
            body.appendChild(separator);

        }


    }
}) 