fetch("/jsonAllTrips/")
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
        const paragraph = document.createElement("p");
        paragraph.setAttribute("id","agency-paragraph"+i);
        child.appendChild(paragraph);
        document.getElementsByTagName("p")[nr].innerHTML="Agentia de turism "+obj.agencyName;
        nr = nr + 1;
        //add the text
        const paragraph2 = document.createElement("p");
        paragraph2.setAttribute("id","title-paragraph"+i);
        child.appendChild(paragraph2);
        document.getElementsByTagName("p")[nr].innerHTML=obj.title;
        nr = nr + 1;
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
        document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
        nr = nr + 1;
        //add the text
        const paragraph4 = document.createElement("p");
        paragraph4.setAttribute("id","price-paragraph"+i);
        child.appendChild(paragraph4);
        document.getElementsByTagName("p")[nr].innerHTML=obj.price+" euro";
        nr = nr + 1;
        //add the text
        const paragraph5 = document.createElement("p");
        paragraph5.setAttribute("id","date-paragraph"+i);
        child.appendChild(paragraph5);
        document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
        //add the photo
        var img1 = document.createElement("img");
        img1.setAttribute("id","img1"+i);        
        img1.src = obj.img1
        child.appendChild(img1)
        //add the photo
        var img2 = document.createElement("img");
        img2.setAttribute("id","img2"+i);        
        img2.src = obj.img2
        child.appendChild(img2)
        //add the photo
        var img3 = document.createElement("img");
        img3.setAttribute("id","img3"+i);        
        img3.src = obj.img3
        child.appendChild(img3)
        nr = nr + 1;
        //add space
        var space = document.createElement("br");
        var separator = document.createElement("hr");
        child.appendChild(space);
        //add a button
        if(obj.same == "yes"){
            var form1 = document.createElement("form");
            form1.setAttribute("id","deleteTrip");
            form1.method="GET";
            form1.action="/deleteTrip/"+obj.id;
            child.appendChild(form1);
            var button = document.createElement("button");
            button.innerHTML = "Sterge excursia";
            button.setAttribute("id","delte-button");
            button.type="submit";
            form1.appendChild(button); 
            child.append(space);

            var form2 = document.createElement("form");
            form2.setAttribute("id","updateTrip");
            form2.method="GET";
            form2.action="/updateTripPage/"+obj.id;
            child.appendChild(form2);
            var button2 = document.createElement("button");
            button2.innerHTML = "Editeaza excursia";
            button2.setAttribute("id","update-button");
            button2.type="submit";
            form2.appendChild(button2); 
            child.append(space)
            }

            var form3 = document.createElement("form");
            form3.setAttribute("id","getWeather");
            form3.method="GET";
            form3.action="/weather/"+obj.city;  
            child.appendChild(form3);
            var button3 = document.createElement("button");
            button3.innerHTML = "Vremea in "+ obj.city;
            button3.type="submit";
            button3.setAttribute("id","weather");
            form3.appendChild(button3); 

        if(obj.is_client == "yes"){
            var form = document.createElement("form");
            form.setAttribute("id","addCart");
            form.method="POST";
            form.action="/addCart/"+obj.id+"/"+obj.clientId;
            child.appendChild(form);
            var button4 = document.createElement("button");
            button4.innerHTML = "Adauga in cos";
            button4.type="submit";
            button4.setAttribute("id","shopping");
            form.appendChild(button4); 
        }
        child.appendChild(separator); 

    }
    console.log(document)
}) 

//search by city
document.getElementById("search").onclick = function search(){
    console.log(22)
    i=0;
    while(document.getElementById("div"+i)){
        document.getElementById("div"+i).hidden = false;
        i++;      
    }

    i=0;
    while(document.getElementById("div"+i)){
        if(document.getElementById("location-paragraph"+i).innerHTML.split("-")[1].toLowerCase().includes(document.getElementById("searchCity").value.toLowerCase())){
            console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
        
    }
    if(document.getElementById("stars1").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 1){
            console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
        }
      }
    
    if(document.getElementById("stars2").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 2){
            console.log(i)
        }
        else
         document.getElementById("div"+i).hidden = true;
         i++;
            
        }
    }
    
    if(document.getElementById("stars3").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 3){
            console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
        }
    }
    
    
    if(document.getElementById("stars4").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 4){
            console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
        }
    }
    
    if(document.getElementById("stars5").checked){
    i=0;
    while(document.getElementById("div"+i)){
        console.log("eu"+parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]));
        console.log(document.getElementById("hotel-paragraph"+i));
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 5){
            console.log(i+"da");
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
            
        }
    }
    
    if(document.getElementById("stars6").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 6){
             console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++;
        }
    }
    
    if(document.getElementById("stars7").checked){
    i=0;
    while(document.getElementById("div"+i)){
        if(parseInt(document.getElementById("hotel-paragraph"+i).innerHTML.split(" ")[3]) == 7){
            console.log(i)
        }
        else
        document.getElementById("div"+i).hidden = true;
        i++; 
        }
    }

    if(document.getElementById("starsall").checked){
        i=0;
        while(document.getElementById("div"+i)){
            if(document.getElementById("location-paragraph"+i).innerHTML.split("-")[1].toLowerCase().includes(document.getElementById("searchCity").value.toLowerCase())){
                document.getElementById("div"+i).hidden = false;
            }
            i++;
            }
        }
    
}


//sort by price
document.getElementById("sortC").onclick = function(){
    fetch("/jsonAllTrips/")
    .then(response => response.json())
    .then(data => {
        var v = Array(data.length);
        for(var i=0; i<data.length; i++){
            v[i] = data[i];
            
        }
        v.sort((a,b) => parseFloat(a.price) - parseFloat(b.price));
        nr=0;
        for(var i = 0; i < v.length; i++) {
            var div = document.getElementById("div"+i);
            div.remove();
        }

        for(var i = 0; i < v.length; i++) {
            var obj = v[i];
            const divAll = document.createElement("div");
            divAll.setAttribute("id","div"+i);
            document.getElementById("body").appendChild(divAll);
            var child = document.getElementById("div"+i);
            //add the text
            const paragraph = document.createElement("p");
            paragraph.setAttribute("id","agency-paragraph"+i);
            child.appendChild(paragraph);
            document.getElementsByTagName("p")[nr].innerHTML="Agentia de turism "+obj.agencyName;
            nr = nr + 1;
            //add the text
            const paragraph2 = document.createElement("p");
            paragraph2.setAttribute("id","title-paragraph"+i);
            child.appendChild(paragraph2);
            document.getElementsByTagName("p")[nr].innerHTML=obj.title;
            nr = nr + 1;
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
            document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
            nr = nr + 1;
            //add the text
            const paragraph4 = document.createElement("p");
            paragraph4.setAttribute("id","price-paragraph"+i);
            child.appendChild(paragraph4);
            document.getElementsByTagName("p")[nr].innerHTML=obj.price+" euro";
            nr = nr + 1;
            //add the text
            const paragraph5 = document.createElement("p");
            paragraph5.setAttribute("id","date-paragraph"+i);
            child.appendChild(paragraph5);
            document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
            //add the photo
            var img1 = document.createElement("img");
            img1.setAttribute("id","img1"+i);        
            img1.src = obj.img1
            child.appendChild(img1)
            //add the photo
            var img2 = document.createElement("img");
            img2.setAttribute("id","img2"+i);        
            img2.src = obj.img2
            child.appendChild(img2)
            //add the photo
            var img3 = document.createElement("img");
            img3.setAttribute("id","img3"+i);        
            img3.src = obj.img3
            child.appendChild(img3)
            nr = nr + 1;
            //add space
            var space = document.createElement("br");
            var separator = document.createElement("hr");
            child.appendChild(space);
            if(obj.same == "yes"){
                var form1 = document.createElement("form");
                form1.setAttribute("id","deleteTrip");
                form1.method="GET";
                form1.action="/deleteTrip/"+obj.id;
                child.appendChild(form1);
                var button = document.createElement("button");
                button.innerHTML = "Sterge excursia";
                button.setAttribute("id","delte-button");
                button.type="submit";
                form1.appendChild(button); 
                child.append(space);
    
                var form2 = document.createElement("form");
                form2.setAttribute("id","updateTrip");
                form2.method="GET";
                form2.action="/updateTripPage/"+obj.id;
                child.appendChild(form2);
                var button2 = document.createElement("button");
                button2.innerHTML = "Editeaza excursia";
                button2.setAttribute("id","update-button");
                button2.type="submit";
                form2.appendChild(button2); 
                child.append(space)
                }
    
                var form3 = document.createElement("form");
                form3.setAttribute("id","getWeather");
                form3.method="GET";
                form3.action="/weather/"+obj.city;  
                child.appendChild(form3);
                var button3 = document.createElement("button");
                button3.innerHTML = "Vremea in "+ obj.city;
                button3.type="submit";
                button3.setAttribute("id","weather");
                form3.appendChild(button3); 
            if(obj.is_client == "yes"){
                var form = document.createElement("form");
                form.setAttribute("id","addCart");
                form.method="POST";
                form.action="/addCart/"+obj.id+"/"+obj.clientId;
                child.appendChild(form);
                var button4 = document.createElement("button");
                button4.innerHTML = "Adauga in cos";
                button4.type="submit";
                button4.setAttribute("id","shopping");
                form.appendChild(button4); 
            }
            child.appendChild(separator); 
    
        }
        document.getElementById("search").click();
    })
}


//sort by price D
document.getElementById("sortD").onclick = function(){
    fetch("/jsonAllTrips/")
    .then(response => response.json())
    .then(data => {
        var v = Array(data.length);
        for(var i=0; i<data.length; i++){
            v[i] = data[i];
            
        }
        v.sort((a,b) => parseFloat(b.price) - parseFloat(a.price));
        nr=0;
        for(var i = 0; i < v.length; i++) {
            var div = document.getElementById("div"+i);
            div.remove();
        }

        for(var i = 0; i < v.length; i++) {
            var obj = v[i];
            const divAll = document.createElement("div");
            divAll.setAttribute("id","div"+i);
            document.getElementById("body").appendChild(divAll);
            var child = document.getElementById("div"+i);
            //add the text
            const paragraph = document.createElement("p");
            paragraph.setAttribute("id","agency-paragraph"+i);
            child.appendChild(paragraph);
            document.getElementsByTagName("p")[nr].innerHTML="Agentia de turism "+obj.agencyName;
            nr = nr + 1;
            //add the text
            const paragraph2 = document.createElement("p");
            paragraph2.setAttribute("id","title-paragraph"+i);
            child.appendChild(paragraph2);
            document.getElementsByTagName("p")[nr].innerHTML=obj.title;
            nr = nr + 1;
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
            document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
            nr = nr + 1;
            //add the text
            const paragraph4 = document.createElement("p");
            paragraph4.setAttribute("id","price-paragraph"+i);
            child.appendChild(paragraph4);
            document.getElementsByTagName("p")[nr].innerHTML=obj.price+" euro";
            nr = nr + 1;
            //add the text
            const paragraph5 = document.createElement("p");
            paragraph5.setAttribute("id","date-paragraph"+i);
            child.appendChild(paragraph5);
            document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
            //add the photo
            var img1 = document.createElement("img");
            img1.setAttribute("id","img1"+i);        
            img1.src = obj.img1
            child.appendChild(img1)
            //add the photo
            var img2 = document.createElement("img");
            img2.setAttribute("id","img2"+i);        
            img2.src = obj.img2
            child.appendChild(img2)
            //add the photo
            var img3 = document.createElement("img");
            img3.setAttribute("id","img3"+i);        
            img3.src = obj.img3
            child.appendChild(img3)
            nr = nr + 1;
            //add space
            var space = document.createElement("br");
            var separator = document.createElement("hr");
            child.appendChild(space);
            if(obj.same == "yes"){
                var form1 = document.createElement("form");
                form1.setAttribute("id","deleteTrip");
                form1.method="GET";
                form1.action="/deleteTrip/"+obj.id;
                child.appendChild(form1);
                var button = document.createElement("button");
                button.innerHTML = "Sterge excursia";
                button.setAttribute("id","delte-button");
                button.type="submit";
                form1.appendChild(button); 
                child.append(space);
    
                var form2 = document.createElement("form");
                form2.setAttribute("id","updateTrip");
                form2.method="GET";
                form2.action="/updateTripPage/"+obj.id;
                child.appendChild(form2);
                var button2 = document.createElement("button");
                button2.innerHTML = "Editeaza excursia";
                button2.setAttribute("id","update-button");
                button2.type="submit";
                form2.appendChild(button2); 
                child.append(space)
                }
    
                var form3 = document.createElement("form");
                form3.setAttribute("id","getWeather");
                form3.method="GET";
                form3.action="/weather/"+obj.city;  
                child.appendChild(form3);
                var button3 = document.createElement("button");
                button3.innerHTML = "Vremea in "+ obj.city;
                button3.type="submit";
                button3.setAttribute("id","weather");
                form3.appendChild(button3); 
            if(obj.is_client == "yes"){
            var form = document.createElement("form");
            form.setAttribute("id","addCart");
            form.method="POST";
            form.action="/addCart/"+obj.id+"/"+obj.clientId;
            child.appendChild(form);
            var button4 = document.createElement("button");
            button4.innerHTML = "Adauga in cos";
            button4.type="submit";
            button4.setAttribute("id","shopping");
            form.appendChild(button4);  
            }
            child.appendChild(separator); 
    
        }
        document.getElementById("search").click();
    })
}


document.getElementById("reset").onclick = function(){
    fetch("/jsonAllTrips/")
    .then(response => response.json())
    .then(data => {
        var nr = 0
        for(var i = 0; i < data.length; i++) {
            var div = document.getElementById("div"+i);
            div.remove();
        }
        for(var i = 0; i < data.length; i++) {
            var obj = data[i];
            const divAll = document.createElement("div");
            divAll.setAttribute("id","div"+i);
            document.getElementById("body").appendChild(divAll);
            var child = document.getElementById("div"+i);
            //add the text
            const paragraph = document.createElement("p");
            paragraph.setAttribute("id","agency-paragraph"+i);
            child.appendChild(paragraph);
            document.getElementsByTagName("p")[nr].innerHTML="Agentia de turism "+obj.agencyName;
            nr = nr + 1;
            //add the text
            const paragraph2 = document.createElement("p");
            paragraph2.setAttribute("id","title-paragraph"+i);
            child.appendChild(paragraph2);
            document.getElementsByTagName("p")[nr].innerHTML=obj.title;
            nr = nr + 1;
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
            document.getElementsByTagName("p")[nr].innerHTML="Hotelul " + obj.hotel +" cu "+obj.stars+" stele";
            nr = nr + 1;
            //add the text
            const paragraph4 = document.createElement("p");
            paragraph4.setAttribute("id","price-paragraph"+i);
            child.appendChild(paragraph4);
            document.getElementsByTagName("p")[nr].innerHTML=obj.price+" euro";
            nr = nr + 1;
            //add the text
            const paragraph5 = document.createElement("p");
            paragraph5.setAttribute("id","date-paragraph"+i);
            child.appendChild(paragraph5);
            document.getElementsByTagName("p")[nr].innerHTML="De la data de "+ obj.date+" timp de "+obj.days+" zile";
            //add the photo
            var img1 = document.createElement("img");
            img1.setAttribute("id","img1"+i);        
            img1.src = obj.img1
            child.appendChild(img1)
            //add the photo
            var img2 = document.createElement("img");
            img2.setAttribute("id","img2"+i);        
            img2.src = obj.img2
            child.appendChild(img2)
            //add the photo
            var img3 = document.createElement("img");
            img3.setAttribute("id","img3"+i);        
            img3.src = obj.img3
            child.appendChild(img3)
            nr = nr + 1;
            //add space
            var space = document.createElement("br");
            var separator = document.createElement("hr");
            child.appendChild(space);
            if(obj.same == "yes"){
                var form1 = document.createElement("form");
                form1.setAttribute("id","deleteTrip");
                form1.method="GET";
                form1.action="/deleteTrip/"+obj.id;
                child.appendChild(form1);
                var button = document.createElement("button");
                button.innerHTML = "Sterge excursia";
                button.setAttribute("id","delte-button");
                button.type="submit";
                form1.appendChild(button); 
                child.append(space);
    
                var form2 = document.createElement("form");
                form2.setAttribute("id","updateTrip");
                form2.method="GET";
                form2.action="/updateTripPage/"+obj.id;
                child.appendChild(form2);
                var button2 = document.createElement("button");
                button2.innerHTML = "Editeaza excursia";
                button2.setAttribute("id","update-button");
                button2.type="submit";
                form2.appendChild(button2); 
                child.append(space)
                }
    
                var form3 = document.createElement("form");
                form3.setAttribute("id","getWeather");
                form3.method="GET";
                form3.action="/weather/"+obj.city;  
                child.appendChild(form3);
                var button3 = document.createElement("button");
                button3.innerHTML = "Vremea in "+ obj.city;
                button3.type="submit";
                button3.setAttribute("id","weather");
                form3.appendChild(button3); 
            if(obj.is_client == "yes"){
                var form = document.createElement("form");
                form.setAttribute("id","addCart");
                form.method="POST";
                form.action="/addCart/"+obj.id+"/"+obj.clientId;
                child.appendChild(form);
                var button4 = document.createElement("button");
                button4.innerHTML = "Adauga in cos";
                button4.type="submit";
                button4.setAttribute("id","shopping");
                form.appendChild(button4); 
            }
            child.appendChild(separator); 
    
        }
        console.log(document)
    }) 
    

}
