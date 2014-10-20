//
//  APITestsTests.swift
//  APITestsTests
//
//  Created by Felix Sun on 10/17/14.
//  Copyright (c) 2014 HyperMusk. All rights reserved.
//

import UIKit
import XCTest
import ApiTests

class APITestsTests: XCTestCase {
    var service = Service()
    
    override func setUp() {
        super.setUp()
        Api.get().baseURL = "http://localhost:9000/api";
        Api.get().verbose = true;

    }
    
    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
        super.tearDown()
    }
    
    func testBlockReturnError() {
        let readyExpectation = expectationWithDescription("ready")
        
        service.permiessionDenied({err in
            XCTAssert(false, "Shouldn't be here")
            readyExpectation.fulfill()
            return
        },
        failure: {err in
            XCTAssertEqual(err.localizedDescription, "permission denied.", "wrong description")
            readyExpectation.fulfill()
        })

        waitForExpectationsWithTimeout(10, { err in
            XCTAssertNil(err, "Error")
        })
    }
    
//    func testPerformanceExample() {
//        // This is an example of a performance test case.
//        self.measureBlock() {
//            // Put the code you want to measure the time of here.
//        }
//    }
    
}
