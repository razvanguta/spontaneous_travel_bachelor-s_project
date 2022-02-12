fetch("/jsonAllAgencies/")
.then(response => response.json())
.then(data => {
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph");
        var body = document.querySelector("body");
        body.appendChild(paragraph);
        document.getElementsByTagName("p")[i].innerHTML="Agentia de turism "+obj.username;
        //add the photo
        var img = document.createElement("img");
        img.setAttribute("id","agency-photo");        
        img.src = obj.profile_image
        body.appendChild(img)
        //add space
        var space = document.createElement("br");
        body.appendChild(space);
        body.appendChild(space);
        //add a button
        var button = document.createElement("a");
        button.innerHTML = "Viziteaza pagina";
        button.href = "/agencyPage/"+obj.username;
        button.setAttribute("id","agency-button");
        body.appendChild(button);
        //add space
        var space = document.createElement("br");
        body.appendChild(space);
        var button2 = document.createElement("a");
        button2.innerHTML = "Vizualizeaza sau adauga recenzii";
        button2.href = "/seeReviews/"+obj.username;
        button2.setAttribute("id","review-button");
        body.appendChild(button2);
        body.appendChild(space);
        if(obj.is_admin == "yes"){
            var button3 = document.createElement("a");
        button3.innerHTML = "Sterge agentia";
        button3.href = "/deleteAgency/"+obj.id;
        button3.setAttribute("id","review-button");
        body.appendChild(button3);
        body.appendChild(space);
        }

    }
}) 