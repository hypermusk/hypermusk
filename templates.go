package main

var Templates = `
{{define "objc/properties_m"}}
{{range . }}{{$f := .ToLanguageField "objc"}}@synthesize {{$f.Name | title}};
{{end}}
- (id) initWithDictionary:(NSDictionary*)dict{
	self = [super init];
	if (!self) {
		return self;
	}
	if (![dict isKindOfClass:[NSDictionary class]]) {
		return self;
	}
{{range . }}{{$f := .ToLanguageField "objc"}}{{ $name := $f.Name | title }}{{if .IsError}}	[self set{{$name}}:[{{$f.Prefix}}{{$f.PkgName | title}} errorWithDictionary:{{$f.SetPropertyObjc}}]];{{else}}{{if $f.Primitive }}{{if $f.IsArray}}	[self set{{$name}}:{{$f.SetPropertyObjc}}];{{else}}	[self set{{$name}}:{{$name | $f.SetPropertyFromObjcDict}}];{{end}}{{else}}{{if $f.IsArray}}
	NSMutableArray * m{{$name}} = [[NSMutableArray alloc] init];
	NSArray * l{{$name}} = [dict valueForKey:@"{{$name}}"];
	if ([l{{$name}} isKindOfClass:[NSArray class]]) {
		for (NSDictionary * d in l{{$name}}) {
			[m{{$name}} addObject: [[{{$f.Prefix}}{{$f.ConstructorType}} alloc] initWithDictionary:d]];
		}
		[self set{{$name}}:m{{$name}}];
	}{{else}}
	id dict{{$name}} = {{$name | $f.SetPropertyFromObjcDict}};
	if ([dict{{$name}} isKindOfClass:[NSDictionary class]]){
		[self set{{$name}}:[[{{$f.Prefix}}{{$f.ConstructorType}} alloc] initWithDictionary:dict{{$name}}]];
	}{{end}}{{end}}{{end}}
{{end}}
	return self;
}

- (NSDictionary*) dictionary {
	NSMutableDictionary * dict = [[NSMutableDictionary alloc] init];
{{range . }}{{$f := .ToLanguageField "objc"}}{{ $name := $f.Name | title }}{{if $f.Primitive }}{{if $f.IsArray}}	[dict setValue:{{$f.GetPropertyObjc}} forKey:@"{{$name}}"];{{else}}	[dict setValue:{{$f.GetPropertyObjc | $f.GetPropertyToObjcDict}} forKey:@"{{$name}}"];{{end}}{{else}}{{if $f.IsArray}}
	NSMutableArray * m{{$name}} = [[NSMutableArray alloc] init];
	for ({{$f.Prefix}}{{$f.Type}} p in {{$name}}) {
		[m{{$name}} addObject:[p dictionary]];
	}
	[dict setValue:m{{$name}} forKey:@"{{$name}}"];{{else}}	[dict setValue:[self.{{$name}} dictionary] forKey:@"{{$name}}"];{{end}}
	{{end}}
{{end}}
	return dict;
}
{{end}}

{{define "objc/properties_h"}}
{{range .}}{{$f := .ToLanguageField "objc"}}@property {{$f.PropertyAnnotation}} {{$f.FullObjcTypeName}} {{.Name | title}};
{{end}}
- (id) initWithDictionary:(NSDictionary*)dict;
- (NSDictionary*) dictionary;
{{end}}

{{define "objc/h"}}// Generated by github.com/hypermusk/hypermusk
// DO NOT EDIT

#import <Foundation/Foundation.h>
{{$pkgName := .Name | title}}
{{$apiprefix := .Prefix}}
@interface {{$apiprefix}}{{$pkgName}} : NSObject
@property (nonatomic, strong) NSString * BaseURL;
@property (nonatomic, assign) BOOL Verbose;
+ ({{$apiprefix}}{{$pkgName}} *) get;

@end


{{range .DataObjects}}{{$do := .}}
// --- {{.Name}} ---
@interface {{$apiprefix}}{{.Name}} : NSObject
{{template "objc/properties_h" $do.Fields}}
@end
{{end}}

// === Interfaces ===
{{range .Interfaces}}{{$interface := .}}
{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
// --- {{$apiprefix}}{{.Name}}Params ---
@interface {{$apiprefix}}{{$interface.Name}}{{.Name}}Params : NSObject
{{template "objc/properties_h" .Params}}
@end

// --- {{$apiprefix}}{{.Name}}Results ---
@interface {{$apiprefix}}{{$interface.Name}}{{.Name}}Results : NSObject
{{template "objc/properties_h" .Results}}
@end
{{end}}{{end}}

@interface {{$apiprefix}}{{.Name}} : NSObject{{with .Constructor}}
{{template "objc/properties_h" .Method.Params}}
{{else}}
- (NSDictionary*) dictionary;
{{end}}
{{range .Methods}}
- {{with .ConstructorForInterface}}({{$apiprefix}}{{.Name}} *){{else}}({{$interface.Name | .ResultsForObjcFunction}}){{end}} {{.ParamsForObjcFunction}};
{{end}}@end
{{end}}
{{end}}

{{define "objc/m"}}// Generated by github.com/hypermusk/hypermusk
// DO NOT EDIT
{{$pkgName := .Name | title}}
{{$apiprefix := .Prefix}}
#import "{{.Name}}.h"

static {{$apiprefix}}{{.Name | title}} * _{{.Name}};
static NSDateFormatter * _dateFormatter;

@implementation {{$apiprefix}}{{.Name | title}} : NSObject
+ ({{$apiprefix}}{{.Name | title}} *) get {
	if(!_{{.Name}}) {
		_{{.Name}} = [[{{$apiprefix}}{{.Name | title}} alloc] init];
	}
	return _{{.Name}};
}

+ (NSDateFormatter *) dateFormatter {
	if(!_dateFormatter) {
		_dateFormatter = [[NSDateFormatter alloc] init];
		[_dateFormatter setDateFormat:@"yyyy-MM-dd'T'HH:mm:ssZZZZ"];
	}
	return _dateFormatter;
}

+ (NSDate *) dateFromString:(NSString *)dateString {
	if(!dateString) {
		return nil;
	}

	NSError *error;
	NSRegularExpression *regexp = [NSRegularExpression regularExpressionWithPattern:@"\\.[0-9]*" options:0 error:&error];
	NSAssert(!error, @"Error in regexp");

	NSRange range = NSMakeRange(0, [dateString length]);
	dateString = [regexp stringByReplacingMatchesInString:dateString options:0 range:range withTemplate:@""];

	NSDate *date;
	[[{{$apiprefix}}{{.Name | title}} dateFormatter] getObjectValue:&date forString:dateString range:nil error:&error];
	if(error) {
		if ([[{{$apiprefix}}{{.Name | title}} get] Verbose]) NSLog(@"Error formatting date %@: %@ (%@)", dateString, [error localizedDescription], error);
		return nil;
	}
	return date;
}

+ (NSString *) stringFromDate:(NSDate *) date {
	if(!date) {
		return nil;
	}
	NSString * dateString = [[{{$apiprefix}}{{.Name | title}} dateFormatter] stringFromDate:date];
	dateString = [[[dateString substringToIndex:(dateString.length - 3)] stringByAppendingString:@":"] stringByAppendingString:[dateString substringFromIndex:(dateString.length - 2)]];
	return dateString;
}

+ (NSDictionary *) request:(NSURL*)url req:(NSDictionary *)req error:(NSError **)error {
	NSMutableURLRequest *httpRequest = [NSMutableURLRequest requestWithURL:url];
	[httpRequest setHTTPMethod:@"POST"];
	[httpRequest setValue:@"application/json;charset=utf-8" forHTTPHeaderField:@"Content-Type"];
	{{$apiprefix}}{{$pkgName}} * _api = [{{$apiprefix}}{{$pkgName}} get];
	NSData *requestBody = [NSJSONSerialization dataWithJSONObject:req options:NSJSONWritingPrettyPrinted error:error];
	if([_api Verbose]) {
		NSLog(@"Request: %@", [NSString stringWithUTF8String:[requestBody bytes]]);
	}
	[httpRequest setHTTPBody:requestBody];
	if(*error != nil) {
		return nil;
	}
	NSURLResponse  *response = nil;
	NSData *returnData = [NSURLConnection sendSynchronousRequest:httpRequest returningResponse:&response error:error];
	if(*error != nil || returnData == nil) {
		return nil;
	}
	if([_api Verbose]) {
		NSLog(@"Response: %@", [NSString stringWithUTF8String:[returnData bytes]]);
	}
	return [NSJSONSerialization JSONObjectWithData:returnData options:NSJSONReadingAllowFragments error:error];
}

+ (NSError *)errorWithDictionary:(NSDictionary *)dict {
	if (![dict isKindOfClass:[NSDictionary class]]) {
		return nil;
	}
	if ([[dict allKeys] count] == 0) {
		return nil;
	}
	NSMutableDictionary *userInfo = [NSMutableDictionary alloc];
	id reason = [dict valueForKey:@"Reason"];
	if ([reason isKindOfClass:[NSDictionary class]]) {
		userInfo = [userInfo initWithDictionary:reason];
	} else {
		userInfo = [userInfo init];
	}
	[userInfo setObject:[dict valueForKey:@"Message"] forKey:NSLocalizedDescriptionKey];

	NSString *code = [dict valueForKey:@"Code"];
	NSNumberFormatter *f = [[NSNumberFormatter alloc] init];
	[f setNumberStyle:NSNumberFormatterDecimalStyle];
	NSNumber *codeNumber = [f numberFromString:code];
	NSInteger intCode = -1;
	if (codeNumber != nil) {
		intCode = [codeNumber integerValue];
	}
	NSError *err = [NSError errorWithDomain:@"{{$pkgName}}Error" code:intCode userInfo:userInfo];
	return err;
}

@end

{{range .DataObjects}}{{$do := .}}
// --- {{.Name}} ---
@implementation {{$apiprefix}}{{.Name}}
{{template "objc/properties_m" $do.Fields}}
@end
{{end}}

// === Interfaces ===

{{range .Interfaces}}{{$interface := .}}
{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
// --- {{$apiprefix}}{{.Name}}Params ---
@implementation {{$apiprefix}}{{$interface.Name}}{{.Name}}Params : NSObject
{{template "objc/properties_m" .Params}}
@end

// --- {{$apiprefix}}{{.Name}}Results ---
@implementation {{$apiprefix}}{{$interface.Name}}{{.Name}}Results : NSObject
{{template "objc/properties_m" .Results}}
@end
{{end}}{{end}}{{end}}

{{range .Interfaces}}{{$interface := .}}
@implementation {{$apiprefix}}{{.Name}} : NSObject
{{with .Constructor}}
{{template "objc/properties_m" .Method.Params}}
{{else}}
- (NSDictionary*) dictionary {
	return [NSDictionary dictionaryWithObjectsAndKeys:nil];
}
{{end}}
{{range .Methods}}{{$method := .}}
// --- {{.Name}} ---
- {{with .ConstructorForInterface}}({{$apiprefix}}{{.Name}} *){{else}}({{$interface.Name | .ResultsForObjcFunction}}){{end}} {{.ParamsForObjcFunction}} {
	{{with .ConstructorForInterface}}
	{{$apiprefix}}{{.Name}} *results = [{{$apiprefix}}{{.Name}} alloc];
	{{range .Constructor.Method.Params}}{{$f := .ToLanguageField "objc"}}[results set{{$f.Name | title}}:{{$f.Name}}];
	{{end}}{{else}}
	{{$apiprefix}}{{$interface.Name}}{{.Name}}Results *results = [{{$apiprefix}}{{$interface.Name}}{{.Name}}Results alloc];
	{{$apiprefix}}{{$interface.Name}}{{.Name}}Params *params = [[{{$apiprefix}}{{$interface.Name}}{{.Name}}Params alloc] init];
	{{range .Params}}{{$f := .ToLanguageField "objc"}}[params set{{$f.Name | title}}:{{$f.Name}}];
	{{end}}
	{{$apiprefix}}{{$pkgName}} * _api = [{{$apiprefix}}{{$pkgName}} get];
	NSURL *url = [NSURL URLWithString:[NSString stringWithFormat:@"%@/{{$interface.Name}}/{{.Name}}.json", [_api BaseURL]]];
	if([_api Verbose]) {
		NSLog(@"Requesting URL: %@", url);
	}
	NSError *error;
	NSDictionary * dict = [{{$apiprefix}}{{$pkgName}} request:url req:[NSDictionary dictionaryWithObjectsAndKeys: [self dictionary], @"This", [params dictionary], @"Params", nil] error:&error];
	if(error != nil) {
		if([_api Verbose]) {
			NSLog(@"Error: %@", error);
		}
		results = [results init];
		[results setErr:error];
		return {{$method.ObjcReturnResultsOrOnlyOne}};
	}
	results = [results initWithDictionary: dict];
	{{end}}
	return {{$method.ObjcReturnResultsOrOnlyOne}};
}
{{end}}@end
{{end}}


{{end}}



{{define "java/properties"}}
{{range . }}{{$f := .ToLanguageField "java"}}	@SerializedName("{{$f.Name | title}}")
	private {{$f.FullJavaTypeName}} _{{$f.Name | snake}};

{{end}}

{{range . }}{{$f := .ToLanguageField "java"}}	public {{$f.FullJavaTypeName}} get{{$f.Name | title }}() {
		return this._{{$f.Name | snake}};
	}
	public void set{{$f.Name | title }}({{$f.FullJavaTypeName}} _{{$f.Name | snake}}) {
		this._{{$f.Name | snake}} = _{{$f.Name | snake}};
	}

{{end}}
{{end}}

{{define "java/packageclass"}}

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.gson.JsonDeserializationContext;
import com.google.gson.JsonDeserializer;
import com.google.gson.JsonElement;
import com.google.gson.JsonParseException;

import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.lang.reflect.Type;
import java.net.HttpURLConnection;
import java.net.URL;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

public class {{.Name | title}} {
	public static {{.Name | title}} _instance;
	public static String dateformat = "yyyy-MM-dd'T'HH:mm:ss.SSSZZZZ";
	private static SimpleDateFormat _formatter = new SimpleDateFormat(dateformat);

	public static {{.Name | title}} get() {
		if (_instance == null) {
			_instance = new {{.Name | title}}();
		}
		return _instance;
	}

	private class DateDeserializer implements JsonDeserializer<Date> {
		public Date deserialize(JsonElement json, Type typeOfT, JsonDeserializationContext context)
				throws JsonParseException {
			String dateString = json.getAsJsonPrimitive().getAsString();

			int len = dateString.length();
			// get rid of ":" in the last timezone 08:00 of 2013-07-16T14:26:41.499+08:00
			if (dateString.charAt(len-3) == ':' && dateString.charAt(len-6) == '+') {
				dateString = dateString.substring(0, len-3) + dateString.substring(len-2, len);
			}

			Date r = null;
			try {
				r = _formatter.parse(dateString);
			} catch (ParseException e) {
			}

			return r;
		}
	}

	public Gson _gson;

	public Gson gson() {
		 if (_gson == null) {
			 GsonBuilder b = new GsonBuilder();
			 b.setDateFormat(dateformat);
			 b.setPrettyPrinting();
			 b.registerTypeAdapter(Date.class, new DateDeserializer());
			 _gson = b.create();
		 }
		return _gson;
	}

	public static class RequestResult {
		private Reader reader;
		private RemoteError err;

		public Reader getReader() {
			return reader;
		}

		public void setReader(Reader reader) {
			this.reader = reader;
		}

		public RemoteError getErr() {
			return err;
		}

		public void setErr(RemoteError err) {
			this.err = err;
		}
	}

	public static RequestResult request(String url, String body) {
		RequestResult r = new RequestResult();
		try {
			URL u = new URL(url);

			HttpURLConnection conn = (HttpURLConnection)u.openConnection();

			conn.setRequestMethod("POST");
			conn.setRequestProperty("Content-Type", "application/json;charset=utf-8");
			conn.setDoInput(true);
			conn.setDoOutput(true);

			DataOutputStream wr = new DataOutputStream(conn.getOutputStream ());
			wr.writeBytes(body);
			wr.flush();
			wr.close();

			InputStream is = conn.getInputStream();

			r.setReader(new InputStreamReader(is));

		} catch (IOException e) {
			RemoteError err = new RemoteError();
			err.setMessage(e.getMessage());
			err.setCode("-1");
			Map reason = new HashMap();
			reason.put("OriginalException", e);
			err.setReason(reason);
			r.setErr(err);
		}

		return r;
	}

	private String baseURL;

	public String getBaseURL() {
		return baseURL;
	}

	public void setBaseURL(String baseURL) {
		this.baseURL = baseURL;
	}

	private boolean verbose;

	public boolean isVerbose() {
		return verbose;
	}

	public void setVerbose(boolean verbose) {
		this.verbose = verbose;
	}

}

{{end}}

{{define "java/remote_error"}}

import java.util.Map;
import com.google.gson.annotations.SerializedName;

public class RemoteError {
	@SerializedName("Code")
	private String _code;
	@SerializedName("Message")
	private String _message;
	@SerializedName("Reason")
	private Map _reason;

	public String getCode() {
		return this._code;
	}
	public void setCode(String _code) {
		this._code = _code;
	}
	public String getMessage() {
		return this._message;
	}
	public void setMessage(String _message) {
		this._message = _message;
	}
	public Map getReason() {
		return this._reason;
	}
	public void setReason(Map _reason) {
		this._reason = _reason;
	}

	@Override
	public String toString() {
		return String.format("%s: %s", this._code, this._message);
	}
}
{{end}}


{{define "java/pkg_object"}}

{{end}}


{{define "java/dataobject"}}
{{if .HasTimeType}}import java.util.Date;{{end}}
{{if .HasArrayType}}import java.util.ArrayList;{{end}}
{{if .HasMapType}}import java.util.Map;{{end}}
import com.google.gson.annotations.SerializedName;
import java.io.Serializable;

public class {{.Name}} implements Serializable {
{{template "java/properties" .Fields}}
}

{{end}}



{{define "java/interface"}}
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Date;
import com.google.gson.annotations.SerializedName;
import com.google.gson.Gson;
import java.io.Serializable;

public class {{.Interface.Name}} implements Serializable {
{{$apiprefix := .Prefix}}{{$interface := .Interface}}{{$pkgName := .PkgName}}{{$interfaceName := .Interface.Name}}
{{with .Interface}}{{with .Constructor}}
{{template "java/properties" .Method.Params}}
{{end}}

{{range .Methods}}{{$method := .}}

	{{with .ConstructorForInterface}}public {{$apiprefix}}{{.Name}} {{else}}
	// --- {{.Name}}Params ---
	public static class {{.Name}}Params {
		{{template "java/properties" .Params}}
	}
	// --- {{.Name}}Results ---
	public static class {{.Name}}Results {
		{{template "java/properties" .Results}}
	}
	// --- {{.Name}} ---
	public {{$interface.Name | .ResultsForJavaFunction}} {{end}}{{.Name | snake}}({{.ParamsForJavaFunction}}) {
	{{with .ConstructorForInterface}}
	{{$apiprefix}}{{.Name}} results = new {{$apiprefix}}{{.Name}}();
	{{range .Constructor.Method.Params}}{{$f := .ToLanguageField "java"}}results.set{{$f.Name | title}}({{$f.Name}});
	{{end}}{{else}}
	{{.Name}}Results results = null;
	{{.Name}}Params params = new {{.Name}}Params();
	{{range .Params}}{{$f := .ToLanguageField "java"}}params.set{{$f.Name | title}}({{$f.Name}});
	{{end}}
	{{$apiprefix}}{{$pkgName}} _api = {{$apiprefix}}{{$pkgName}}.get();

	Gson gson = _api.gson();

	Map m = new HashMap();
	m.put("Params", params);
	m.put("This", this);

	String url = String.format("%s/{{$interfaceName}}/{{$method.Name}}.json", _api.getBaseURL());
	String requestBody = gson.toJson(m);
	Qortexapi.RequestResult r = Qortexapi.request(url, requestBody);

	if (r.getErr() != null) {
		results = new {{.Name}}Results();
		results.setErr(r.getErr());
		return {{$method.JavaReturnResultsOrOnlyOne}};
	}

	results = gson.fromJson(r.getReader(), {{$method.Name}}Results.class);
	try {
		r.getReader().close();
	} catch (IOException e) {
	}

	{{end}}

	return {{$method.JavaReturnResultsOrOnlyOne}};
	}
{{end}}
{{end}}
}

{{end}}



{{define "httpserver"}}// Generated by github.com/hypermusk/hypermusk
// DO NOT EDIT
{{$pkg := .}}
package {{.Name}}httpimpl

import ({{range .ServerImports}}
	"{{.}}"{{end}}
)

var _ = time.Sunday

type CodeError interface {
	Code() string
}

type SerializableError struct {
	Code    string
	Message string
	Reason  error
}

func (s *SerializableError) Error() string {
	return s.Message
}

func NewError(err error) (r error) {
	se := &SerializableError{Message:err.Error()}
	ce, yes := err.(CodeError)
	if yes {
		se.Code = ce.Code()
	}
	se.Reason = err
	r = se
	return
}

func AddToMux(prefix string, mux *http.ServeMux) {
	{{range .Interfaces}}{{$interface := .}}{{range .Methods}}{{if .ConstructorForInterface}}{{else}}
	mux.HandleFunc(prefix+"/{{$interface.Name}}/{{.Name}}.json", {{$interface.Name}}_{{.Name}}){{end}}{{end}}{{end}}
	return
}
{{range .Interfaces}}{{$interface := .}}
{{with .Constructor}}{{else}}
var {{$interface.Name | downcase}} {{$pkg.Name}}.{{$interface.Name}} = {{$pkg.ImplPkg | dotlastname}}.Default{{$interface.Name}}{{end}}

type {{$interface.Name}}Data struct {
{{with .Constructor}}{{range .Method.Params}}	{{.Name | title}} {{.FullGoTypeName}}
{{end}}{{end}}}

{{range .Methods}}{{if .ConstructorForInterface}}{{else}}{{$method := .}}
type {{$interface.Name}}_{{$method.Name}}_Params struct {
{{with $interface.Constructor}}	This   {{$interface.Name}}Data
{{end}}	Params struct {
{{range .Params}}		{{.Name | title}} {{.FullGoTypeName}}
{{end}}	}
}

type {{$interface.Name}}_{{$method.Name}}_Results struct {
{{range .Results}}	{{.Name | title}} {{.FullGoTypeName}}
{{end}}
}

func {{$interface.Name}}_{{$method.Name}}(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	var p {{$interface.Name}}_{{$method.Name}}_Params
	if r.Body == nil {
		panic("no body")
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	var result {{$interface.Name}}_{{$method.Name}}_Results
	enc := json.NewEncoder(w)
	if err != nil {
		result.Err = NewError(err)
		enc.Encode(result)
		return
	}
{{if $interface.Constructor}}
	s, err := {{$interface.Constructor.FromInterface.Name | downcase }}.{{$interface.Constructor.Method.Name}}({{$interface.Constructor.Method.ParamsForGoServerConstructorFunction}})
{{else}}
	s := {{$interface.Name | downcase }}
{{end}}
	if err != nil {
		result.Err = NewError(err)
		enc.Encode(result)
		return
	}
	{{$method.ResultsForGoServerFunction "result"}} = s.{{$method.Name}}({{$method.ParamsForGoServerFunction}})
	if result.Err != nil {
		result.Err = NewError(result.Err)
	}
	err = enc.Encode(result)
	if err != nil {
		panic(err)
	}
	return
}
{{end}}{{end}}

{{end}}





{{end}}

{{define "javascript/interfaces"}}// Generated by github.com/hypermusk/hypermusk
// DO NOT EDIT

(function(api, $, undefined ) {
	api.rpc = function(endpoint, input, callback) {
		var methodUrl = api.baseurl + endpoint;
		var message = JSON.stringify(input);
		var req = $.ajax({
			type: "POST",
			url: methodUrl,
			contentType:"application/json; charset=utf-8",
			dataType:"json",
			processData: false,
			data: message
		});
		req.done(function(data, textStatus, jqXHR) {
			callback(data);
		});
	};
})( window.{{.Name}} = window.{{.Name}} || {}, jQuery);



(function( api, undefined ) {
{{range .Interfaces}}{{ $interfaceName := .Name}}
	api.{{$interfaceName}} = function() {};
{{range .Methods}}{{$method := .}}{{if .ConstructorForInterface}}
	api.{{$interfaceName}}.prototype.{{.Name}} = function({{$method.ParamsForJavascriptFunction}}) {
		var r = new api.{{.ConstructorForInterface.Name}}(){{range .Params}};
		r.{{.Name | title}} = {{.Name}}{{end}};
		return r;
	}
{{else}}
	api.{{$interfaceName}}.prototype.{{.Name}} = function({{$method.ParamsForJavascriptFunction}}{{if $method.ParamsForJavascriptFunction}}, {{end}}callback) {
		api.rpc("/{{$interfaceName}}/{{.Name}}.json", {"This": this, "Params": {{$method.ParamsForJson}}}, function(data){
			callback({{$method.ResultsForJavascriptFunction "data"}})
		});
		return;
	}
{{end}}{{end}}{{end}}

}( window.{{.Name}} = window.{{.Name}} || {} ));

{{end}}

`