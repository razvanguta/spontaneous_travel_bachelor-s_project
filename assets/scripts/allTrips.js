fetch("/jsonAllTrips/")
.then(response => response.json())
.then(data => {
    var nr = 0
    for(var i = 0; i < data.length; i++) {
        console.log(1)
        var obj = data[i];
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph);
        document.getElementsByTagName("p")[nr].innerHTML="Agentia de turism "+obj.agencyName;
        nr = nr + 1;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","title-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph2);
        document.getElementsByTagName("p")[nr].innerHTML=obj.title;
        nr = nr + 1;
        //add the text
        const paragraph3 = document.createElement("p");
        paragraph3.setAttribute("id","desc-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph3);
        document.getElementsByTagName("p")[nr].innerHTML=obj.description;
        nr = nr + 1;
        //add the text
        const paragraph6 = document.createElement("p");
        paragraph6.setAttribute("id","hotel-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph6);
        document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
        nr = nr + 1;
        //add the text
        const paragraph4 = document.createElement("p");
        paragraph4.setAttribute("id","price-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph4);
        document.getElementsByTagName("p")[nr].innerHTML=obj.price+" euro";
        nr = nr + 1;
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","title-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph5);
        document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
        //add the photo
        var img1 = document.createElement("img");
        img1.setAttribute("id","agency-photo");        
        img1.src = obj.img1
        body.appendChild(img1)
        //add the photo
        var img2 = document.createElement("img");
        img2.setAttribute("id","agency-photo");        
        img2.src = obj.img2
        body.appendChild(img2)
        //add the photo
        var img3 = document.createElement("img");
        img3.setAttribute("id","agency-photo");        
        img3.src = obj.img3
        body.appendChild(img3)
        nr = nr + 1;
        //add space
        var space = document.createElement("br");
        var separator = document.createElement("hr");
        body.appendChild(space);
        //add a button
        if(obj.same == "yes"){
        var button = document.createElement("a");
        button.innerHTML = "Sterge excursia";
        button.setAttribute("id","delete-button");
        button.href = "/deleteTrip/"+obj.id;   
        body.appendChild(button);

        body.append(space)

        var button2 = document.createElement("a");
        button2.innerHTML = "Editeaza excursia";
        button2.setAttribute("id","update-button");
        button2.href = "/UpdateTripPage/"+obj.id;  
        body.appendChild(button2); 
        }
        body.appendChild(separator);

    }
}) 