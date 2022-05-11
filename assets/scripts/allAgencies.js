fetch("/jsonAllAgencies/")
.then(response => response.json())
.then(data => {
    for(var i = 0; i < data.length; i++) {
        var obj = data[i];
        //add the text
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph");
        var body = document.getElementById("body");
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
        
        //add button
        var form = document.createElement("form");
        form.setAttribute("id","agencyPage");
        form.method="get";
        form.action="/agencyPage/"+obj.username;
        var button = document.createElement("button");
        button.innerHTML = "Viziteaza pagina";
        button.type = "submit";
        button.setAttribute("id","agency-button");
        form.appendChild(button);
        body.appendChild(form);
        //add space
        var space = document.createElement("br");
        var form2 = document.createElement("form");
        form2.setAttribute("id","seeReviews");
        form2.method="get";
        form2.action="/seeReviews/"+obj.username;
        body.appendChild(space);
        var button2 = document.createElement("button");
        button2.innerHTML = "Vizualizeaza sau adauga recenzii";
        button2.setAttribute("id","review-button");
        form2.appendChild(button2);
        body.appendChild(form2);


        if(obj.is_admin == "yes"){
        var form3 = document.createElement("form");
        form3.setAttribute("id","deleteAgency");
        form3.method="get";
        form3.action="/deleteAgency/"+obj.id;
        var button3 = document.createElement("button");
        button3.innerHTML = "Sterge agentia";
        button3.setAttribute("id","review-button");
        form3.appendChild(button3);
        body.appendChild(form3);
        body.appendChild(space);
        }

    }
}) 

document.getElementById("body").createElement('footer').appendChild()