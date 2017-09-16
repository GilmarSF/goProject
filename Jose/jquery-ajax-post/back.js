$( document ).ready(function() {
	
});


    
function Save() {		
debugger
	var apiurl = 'http://localhost:8080/categorias/';
    var a = new XMLHttpRequest();
    var datas = {
            "categoria":$('#categoria').val(), "localidade":$('#localidade').val()
        };
		
    a.open('POST',apiurl,true),
    $.ajax({
		url: apiurl,
		type: "POST",
		ContentType: "application/json; charset=utf-8",
		data: JSON.stringify(datas),
		crossOrigin: true,
		dataType: "json",
		cache: false,
		complete: function (res) {
			alert("Data Added Successfully");
		},
		error: function (xhr) {
			alert("Error");
		}
	})

    }